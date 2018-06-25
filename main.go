package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"RPN/model"
	"encoding/json"
	"RPN/dao"
	"gopkg.in/mgo.v2/bson"
	"github.com/dgrijalva/jwt-go"
	"time"
	"RPN/config"
)
var userDao = dao.UserDAO{}


func main() {
	r := mux.NewRouter()
	r.Handle("/register", http.HandlerFunc(signupHandler)).Methods("POST")
	r.Handle("/login", http.HandlerFunc(loginHandler)).Methods("POST")
	//r.Handle("/findUser", jwtMiddleware.Handler( http.HandlerFunc(FindUserHandler))).Methods("POST")
	r.Handle("/findAll", http.HandlerFunc(getAllUsersHandler)).Methods("GET")
	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received new sign up request")
	defer r.Body.Close()
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.ID = bson.NewObjectId()
	//
	if err := userDao.AddUser(user); err != nil {
		log.Println("DB insert error")
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondWithJson(w, http.StatusCreated, user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received new login request")
	defer r.Body.Close()
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if userDao.CheckUser(user.Username, user.Password) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, _ := token.SignedString(config.MySigningKey)
		w.Write([]byte(tokenString))
		log.Println("tokenString: " + tokenString)
		log.Println("Login successfully")
	} else {
		log.Println("Invalid username or password.")
		respondWithError(w, http.StatusForbidden, "Invalid username or password")
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func findUserHandler(w http.ResponseWriter, r *http.Request) {

}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := userDao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, users)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.WriteHeader(code)
	w.Write(response)
}
