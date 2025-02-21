package main

import (
	"authApi/handler"
	"authApi/middleware"
	"authApi/model"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var database map[string]model.DB
	database = make(map[string]model.DB)

	var revocationDb map[float64]bool
	revocationDb = make(map[float64]bool)

	r := mux.NewRouter()

	h := handler.New(database, revocationDb)
	a := middleware.New(revocationDb)

	r.HandleFunc("/signUp", h.SignUp).Methods("POST")
	r.HandleFunc("/signIn", h.SignIn).Methods("POST")
	r.Handle("/get", a.AuthenticationMiddleware(http.HandlerFunc(h.GetOperation))).Methods("GET")
	r.HandleFunc("/refreshToken", h.RefreshToken).Methods("GET")
	r.HandleFunc("/revokeToken", h.RevokeToken).Methods("GET")

	fmt.Println("Listening on port: 8000...")

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		fmt.Println("Unable to start server:", err)
	}
}
