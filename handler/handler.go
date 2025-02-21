package handler

import (
	"authApi/model"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	secretKey        = []byte("Ariqt")
	refreshSecretKey = []byte("AriqtRefresh")
)

type handler struct {
	Database           *map[string]model.DB
	RevocationDatabase map[float64]bool
}

func New(db map[string]model.DB, revocationDb map[float64]bool) handler {
	return handler{Database: &db, RevocationDatabase: revocationDb}
}

func (h handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userCreds model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("invalid request body:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &userCreds)
	if err != nil {
		fmt.Errorf("invalid request body:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(userCreds.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Errorf("please choose proper password:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	(*h.Database)[userCreds.UserName] = model.DB{Password: pwd, IsUserValid: true}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Sign Up Successful!"))
}

func (h handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var userCreds model.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("invalid request body:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &userCreds)
	if err != nil {
		fmt.Errorf("invalid request body:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pwd := (*h.Database)[userCreds.UserName]

	err = bcrypt.CompareHashAndPassword(pwd.Password, []byte(userCreds.Password))
	if err != nil {
		fmt.Errorf("wrong password:%s", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("wrong password"))
		return
	}

	tkn := generateJWT(userCreds.UserName)
	if tkn == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshTkn := generateRefreshToken(userCreds.UserName)
	if refreshTkn == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := model.SignInResponse{Token: tkn, RefreshToken: refreshTkn}

	token, err := json.Marshal(resp)
	if err != nil {
		fmt.Errorf("internal server error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(token)
}

func generateJWT(userName string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userName
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["jti"] = rand.Float64()

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Errorf("internal server error:%s", err.Error())
		return ""
	}

	return signedToken
}

func generateRefreshToken(userName string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userName
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	signedToken, err := token.SignedString(refreshSecretKey)
	if err != nil {
		fmt.Errorf("internal server error:%s", err.Error())
		return ""
	}

	return signedToken
}

func (h handler) GetOperation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authenticated successfully"))
}

func (h handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Refresh-Token")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter refresh token"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter a Bearer Refresh Token"))
		return
	}

	token, err := jwt.ParseWithClaims(parts[1], &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return refreshSecretKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter a valid refresh token"))
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter a valid refresh token"))
		return
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Please enter a valid refresh token"))
		return
	}

	if claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Refresh Token expired"))
		return
	}

	usrName, ok := (*claims)["username"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter a valid refresh token"))
		return
	}

	resp := model.SignInResponse{Token: generateJWT(usrName)}

	tokenResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Errorf("internal server error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(tokenResp)
}

func (h handler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter token"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter a Bearer token"))
		return
	}

	token, err := jwt.ParseWithClaims(parts[1], &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("Ariqt"), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter a valid token"))
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter a valid token"))
		return
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Please enter a valid token"))
		return
	}

	if claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Token already expired"))
		return
	}

	jti, ok := (*claims)["jti"].(float64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter a valid token"))
		return
	}

	h.RevocationDatabase[jti] = true

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Revocation successful!"))
}
