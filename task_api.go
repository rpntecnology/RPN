package main

import (
	"RPN/dao"
	"net/http"
	"log"
	"RPN/model"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"github.com/dgrijalva/jwt-go"

	"strconv"
)

var taskDao = dao.TaskDAO{}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) < AUTH_TO_MANAGE_TASK {
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
	log.Println("Received new add image request")
	defer r.Body.Close()
	var imageSlot model.ImageSlot
	if err := json.NewDecoder(r.Body).Decode(&imageSlot); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	imageSlot.ImageID = bson.NewObjectId()
	if err := taskDao.AddImage(imageSlot); err != nil {
		log.Println("DB insert error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondWithJson(w, http.StatusCreated, imageSlot)
}


func FindImgURLByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Find images by category request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	taskId, _ := r.URL.Query().Get("task_id"), 64
	category, _ := r.URL.Query().Get("category"), 64
	log.Println(taskId)
	log.Println(category)

	err, task := taskDao.FindById(bson.ObjectIdHex(taskId))
	if err != nil {
		log.Println("Error in finding task")
		respondWithError(w, http.StatusInternalServerError, "Error in finding task")
		return
	}
	//log.Println(task.Image)

	images := task.Image
	var urls []string
	for _, image := range images {
		if image.Category == category {
			urls = append(urls, image.URL)
		}
	}
	log.Println(urls)
	respondWithJson(w, http.StatusOK, urls)
}

func checkAuth(r *http.Request) int {
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	authority := claims.(jwt.MapClaims)["authority"]
	auth, _ := strconv.Atoi(authority.(string))
	return auth
}