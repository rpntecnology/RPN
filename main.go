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
	r.Handle("/removeUser", jwtMiddleware.Handler(http.HandlerFunc(RemoveUserHandler))).Methods("POST")
	r.Handle("/addTask", jwtMiddleware.Handler(http.HandlerFunc(AddTaskHandler))).Methods("POST")
	r.Handle("/addImage", jwtMiddleware.Handler(http.HandlerFunc(AddImageHandler))).Methods("POST")
	r.Handle("/deleteImageFromTask", jwtMiddleware.Handler(http.HandlerFunc(DeleteImageHandler))).Methods("POST")
	r.Handle("/deleteTask", jwtMiddleware.Handler(http.HandlerFunc(DeleteTaskHandler))).Methods("POST")
	r.Handle("/ChangeTaskUser", jwtMiddleware.Handler(http.HandlerFunc(ChangeContractorHandler))).Methods("POST")
	//r.Handle("/findImageByCategory", http.HandlerFunc(FindImgURLByCategoryHandler)).Methods("GET")
	r.Handle("/findAll", http.HandlerFunc(GetAllUsersHandler)).Methods("GET")
	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}




