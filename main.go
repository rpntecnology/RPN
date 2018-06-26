package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"RPN/config"
)

func main() {
	r := mux.NewRouter()
	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return config.MySigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	r.Handle("/register", http.HandlerFunc(SignupHandler)).Methods("POST")
	r.Handle("/login", http.HandlerFunc(LoginHandler)).Methods("POST")
	r.Handle("/userProfile", jwtMiddleware.Handler(http.HandlerFunc(UserProfileHandler))).Methods("GET")
	r.Handle("/addTask", jwtMiddleware.Handler(http.HandlerFunc(AddTaskHandler))).Methods("POST")
	r.Handle("/findAll", http.HandlerFunc(GetAllUsersHandler)).Methods("GET")
	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}




