package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"RPN/model"
	"encoding/json"
	"RPN/dao"
	"gopkg.in/mgo.v2/bson"
)
var userDao = dao.UserDAO{}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", http.HandlerFunc(signupHandler)).Methods("POST")
	r.HandleFunc("/findAll", getAllUsersHandler).Methods("GET")
	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received new sign up request!")
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
