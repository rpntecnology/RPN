package main

import (
	"RPN/dao"
	"net/http"
	"log"
	"RPN/model"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"github.com/dgrijalva/jwt-go"

	"strings"
)

var taskDao = dao.TaskDAO{}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		log.Println("Only admin can add tasks!")
		respondWithError(w, http.StatusForbidden, "Only admin can add tasks!")
		return
	}

	log.Println("Received new add task request")
	defer r.Body.Close()
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	task.TaskID = bson.NewObjectId()
	//
	if err := taskDao.AddTask(task); err != nil {
		log.Println("DB insert error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondWithJson(w, http.StatusCreated, task)
}

func AddImageHandler(w http.ResponseWriter, r *http.Request) {

}

func isAdmin(r *http.Request) bool {
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]
	return strings.Compare(userDao.FindAuthority(username.(string)), "admin") == 0
}