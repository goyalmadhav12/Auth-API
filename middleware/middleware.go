package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type auth struct {
	RevocationDatabase map[float64]bool
}

func New(revocationDb map[float64]bool) auth {
	return auth{RevocationDatabase: revocationDb}
}

func (a auth) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			w.Write([]byte("Token expired"))
			return
		}

		jti, ok := (*claims)["jti"].(float64)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Please enter a valid token"))
			return
		}

		if a.RevocationDatabase[jti] == true {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token already revoked"))
			return
		}

		fmt.Println("Authenticated successfully")
		next.ServeHTTP(w, r)
	})
}
