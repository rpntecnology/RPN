package main

import (
	"github.com/gorilla/mux"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"RPN/config"
	"net/http"
	"log"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return config.MySigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	r.Handle(config.API_PREFIX + "/register", http.HandlerFunc(SignupHandler)).Methods("POST")
	r.Handle(config.API_PREFIX + "/login", http.HandlerFunc(LoginHandler)).Methods("POST")
	r.Handle(config.API_PREFIX + "/userProfile", jwtMiddleware.Handler(http.HandlerFunc(UserProfileHandler))).Methods("GET")
	r.Handle(config.API_PREFIX + "/removeUser", jwtMiddleware.Handler(http.HandlerFunc(RemoveUserHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/addTask", jwtMiddleware.Handler(http.HandlerFunc(AddTaskHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/initTask", jwtMiddleware.Handler(http.HandlerFunc(InitTaskHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/updateTask", jwtMiddleware.Handler(http.HandlerFunc(UpdateTaskHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/addImage", jwtMiddleware.Handler(http.HandlerFunc(AddImageHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/deleteImageFromTask", jwtMiddleware.Handler(http.HandlerFunc(DeleteImageHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/deleteTask", jwtMiddleware.Handler(http.HandlerFunc(DeleteTaskHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/addCategory", jwtMiddleware.Handler(http.HandlerFunc(AddCategoryHandler))).Methods("POST")
	//r.Handle("/addItem", jwtMiddleware.Handler(http.HandlerFunc(AddItemHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/addTaskToUser", jwtMiddleware.Handler(http.HandlerFunc(AddTaskToUserHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/changeTaskUser", jwtMiddleware.Handler(http.HandlerFunc(ChangeContractorHandler))).Methods("POST")
	r.Handle(config.API_PREFIX + "/findImg", http.HandlerFunc(FindImgURLHandler)).Methods("GET")
	r.Handle(config.API_PREFIX + "/findAll", http.HandlerFunc(GetAllUsersHandler)).Methods("GET")



	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	handler := c.Handler(r)

	http.Handle("/", r)

	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}






