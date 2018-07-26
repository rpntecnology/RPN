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
	"strconv"
)

var userDao = dao.UserDAO{}
const (
	AUTH_TO_DELETE = 2
	AUTH_TO_MANAGE_TASK = 1
)

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
	w.Header().Set("Content-Type", "application/json")
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

	err, userDb := userDao.FindUser(user.Username)
	if err != nil {
		log.Println(err.Error())
		log.Println("Error in finding user")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if userDb.Password == user.Password {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = userDb.Username
		claims["authority"] = userDb.Authority
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, _ := token.SignedString(config.MySigningKey)
		type Response struct {
			Token 		string
			Authority 	string
		}
		var response Response
		response.Token = tokenString
		response.Authority = userDb.Authority

		//w.Write([]byte(tokenString))
		respondWithJson(w, http.StatusOK, response)
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
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]
	log.Println(username)
	taskIds := getUsersTaskIds(username.(string))
	tasks := getUsersTasks(taskIds)


	respondWithJson(w, http.StatusOK, tasks)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received get all users request")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	authority := claims.(jwt.MapClaims)["authority"]
	auth, _ := strconv.Atoi(authority.(string))

	if auth < AUTH_TO_DELETE {
		respondWithError(w, http.StatusUnauthorized, "You don not have the authority to view all users")
		return
	}

	users, err := userDao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Print(users)
	respondWithJson(w, http.StatusOK, users)
}

func RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received delete user request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	authority := claims.(jwt.MapClaims)["authority"]
	auth, _ := strconv.Atoi(authority.(string))

	if auth < AUTH_TO_DELETE {
		respondWithError(w, http.StatusUnauthorized, "You don not have the authority to remove user")
		return
	}
	username, _ := r.URL.Query().Get("username"), 64
	log.Println(username)
	err := userDao.DeleteUser(username)
	if err != nil {
		log.Println("Error in deleteing user: " + username)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, "done")
}

func AddTaskToUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received add task to user request")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	if CheckAuth(r) < AUTH_TO_MANAGE_TASK {
		log.Println("No authority to manage task")
		respondWithError(w, http.StatusForbidden, "No authority to manage task")
		return
	}

	userToRemove := r.FormValue("userToRemove")
	userToAdd := r.FormValue("userToAdd")
	taskId := r.FormValue("task_id")
	stage := r.FormValue("stage")

	log.Println("debug")
	log.Println(userToRemove)
	log.Println(userToRemove)
	log.Println(taskId)
	log.Print(stage)

	if _, err := strconv.Atoi(stage); err != nil {
		log.Print("stage is not a number")
		respondWithError(w, http.StatusForbidden, "stage is not a number")
		return
	}

	if userToRemove == "" {
		userToRemove = userToAdd
	}
	err1 := userDao.AssignTaskToUser(userToRemove, userToAdd, bson.ObjectIdHex(taskId))

	if err1 != nil {
		log.Println("Error in assigning tasks, err: " + err1.Error())
		respondWithError(w, http.StatusForbidden, "Error in assigning task to user")
		return
	}

	err2 := taskDao.AssignUserToTask(bson.ObjectIdHex(taskId), userToAdd, stage)
	if err2 != nil {
		log.Println("Error in assigning users, err: " + err2.Error())
		respondWithError(w, http.StatusForbidden, "Error in assigning user to task")
		return
	}
	respondWithJson(w, http.StatusOK, "done")
}



func getUsersTaskIds(username string) []bson.ObjectId {
	err, user := userDao.FindUser(username)
	if err != nil {
		log.Println(err.Error())
	}
	return user.TaskIds
}

func getUsersTasks(taskIds []bson.ObjectId) []model.Task{
	var tasks []model.Task

	for _, taskId := range taskIds {
		err, task := taskDao.FindByIdPreview(taskId)
		if err != nil {
			log.Println(err.Error())
		}
		tasks = append(tasks, task)
	}
	return tasks
}

