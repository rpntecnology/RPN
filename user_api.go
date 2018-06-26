package main

import (
	"net/http"
	"log"
	"RPN/model"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"github.com/dgrijalva/jwt-go"
	"time"
	"RPN/config"
	"RPN/dao"
)

var userDao = dao.UserDAO{}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondWithJson(w, http.StatusCreated, user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Println("username: " + user.Username + "password: " + user.Password)
		respondWithError(w, http.StatusForbidden, "Invalid username or password")
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received new profile request")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	user := r.Context().Value("user")
	//log.Println(user)
	claims := user.(*jwt.Token).Claims
	//log.Println(claims)
	username := claims.(jwt.MapClaims)["username"]

	err, profile := userDao.FindUser(username.(string))
	if err != nil {
		respondWithError(w, http.StatusNoContent, err.Error())
	}
	respondWithJson(w, http.StatusOK, profile)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := userDao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, users)
}
