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
	"context"
	"io"
	"cloud.google.com/go/storage"
	"path/filepath"
	"RPN/config"
)

var taskDao = dao.TaskDAO{}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) < AUTH_TO_MANAGE_TASK {
		log.Println("No authority to add task")
		respondWithError(w, http.StatusForbidden, "No authority to add task")
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
	if checkAuth(r) < 0 {
		log.Println("No authority to post images")
		respondWithError(w, http.StatusForbidden, "No authority to post images")
		return
	}
	log.Println("Received new add image request")
	defer r.Body.Close()
	var imageSlot model.ImageSlot
	//if err := json.NewDecoder(r.Body).Decode(&imageSlot); err != nil {
	//	respondWithError(w, http.StatusInternalServerError, err.Error())
	//	return
	//}


	// 32 << 20 is the maxMemory param for ParseMultipartForm, equals to 32MB (1MB = 1024 * 1024 bytes = 2^20 bytes)
	// After you call ParseMultipartForm, the file will be saved in the server memory with maxMemory size.
	// If the file size is larger than maxMemory, the rest of the data will be saved in a system temporary file.
	r.ParseMultipartForm(32 << 20)
	taskId, _ := r.FormValue("task_id"), 64
	name, _ := r.FormValue("name"), 64
	cate, _ := r.FormValue("cate"), 64
	itemId, _ := r.FormValue("itemId"), 64
	status, _ := r.FormValue("status"), 64
	log.Println("itemId: " + itemId)
	log.Println("taskId: " + taskId)
	log.Println("cate: " + cate)
	log.Println("status" + status)
	imageSlot.ImageID = bson.NewObjectId()
	imageSlot.TaskID = bson.ObjectIdHex(taskId)
	imageSlot.Name = name
	imageSlot.Cate = cate
	imageSlot.ItemId = itemId
	imageSlot.Status = status
	//log.Println(imageSlot.TaskID)
	file, _, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Image is not available")
		log.Println("Image is not available, err: " + err.Error())
		return
	}
	defer file.Close()

	// read image's extension
	img, header, _ := r.FormFile("image")
	defer img.Close()
	suffix := filepath.Ext(header.Filename)

	if _, ok := config.MediaTypes[suffix]; ok {
		imageSlot.Format = suffix
	} else {
		imageSlot.Format = "unknown"
	}

	ctx := context.Background()
	_, attrs, err := saveToGCS(ctx, file, config.BUCKET_NAME, imageSlot.ImageID.String())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,"GCS is not setup")
		log.Println("GCS is not setup, err: " + err.Error())
		return
	}

	imageSlot.Src = attrs.MediaLink


	if err := taskDao.AddPrevImage(imageSlot); err != nil {
		log.Println("DB insert error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondWithJson(w, http.StatusCreated, imageSlot)
}


//func FindImgURLByCategoryHandler(w http.ResponseWriter, r *http.Request) {
//	log.Println("Received Find images by category request")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Content-Type", "application/json")
//	taskId, _ := r.URL.Query().Get("task_id"), 64
//	category, _ := r.URL.Query().Get("category"), 64
//	log.Println(taskId)
//	log.Println(category)
//
//	err, task := taskDao.FindById(bson.ObjectIdHex(taskId))
//	if err != nil {
//		log.Println("Error in finding task")
//		respondWithError(w, http.StatusInternalServerError, "Error in finding task")
//		return
//	}
//	//log.Println(task.Image)
//
//	images := task.Image
//	var urls []string
//	for _, image := range images {
//		if image.Category == category {
//			urls = append(urls, image.Src)
//		}
//	}
//	log.Println(urls)
//	respondWithJson(w, http.StatusOK, urls)
//}


func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) < 0 {
		respondWithError(w, http.StatusInternalServerError, "No authority to delete image")
		log.Println("No authority to delete image")
		return
	}
	log.Println("Received delete image request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	taskId, _ := r.URL.Query().Get("task_id"), 64
	imageId, _ := r.URL.Query().Get("image_id"), 64
	log.Println(taskId)
	log.Println(imageId)

	//var image []model.ImageSlot
	if err := taskDao.DeleteImageByImageID(bson.ObjectIdHex(taskId), bson.ObjectIdHex(imageId)); err != nil {
		log.Println("taskID: "+taskId + " imageID: " + imageId)
		log.Println("DB find error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	log.Println("deleted image successfully")
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) < AUTH_TO_DELETE {
		respondWithError(w, http.StatusInternalServerError, "No authority to delete task")
		log.Println("No authority to delete task")
		return
	}
	log.Println("Received delete task request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	taskId, _ := r.URL.Query().Get("task_id"), 64
	log.Println(taskId)

	if err := taskDao.DeleteTaskByTaskID(bson.ObjectIdHex(taskId)); err != nil {
		log.Println("taskID: "+taskId )
		log.Println("DB find error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	log.Println("deleted task successfully")
}

func ChangeContractorHandler(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) < AUTH_TO_DELETE {
		respondWithError(w, http.StatusInternalServerError, "No authority to assign work")
		log.Println("No authority to assign work")
		return
	}
	log.Println("Received change task's contractor request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	taskId, _ := r.URL.Query().Get("task_id"), 64
	log.Println(taskId)
	user, _ := r.URL.Query().Get("username"), 64
	log.Println(user)

	if err := taskDao.AssignTaskToAnotherUser(bson.ObjectIdHex(taskId), user); err != nil {
		log.Println("taskID: "+taskId +" user: " + user)
		log.Println("DB find error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	log.Println("changed task's contractor successfully")
}


func checkAuth(r *http.Request) int {
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	authority := claims.(jwt.MapClaims)["authority"]
	auth, _ := strconv.Atoi(authority.(string))
	return auth
}

func saveToGCS(ctx context.Context, r io.Reader, bucketName, name string) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	// check if the bucket exists
	if _, err = bucket.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bucket.Object(name)
	wc := obj.NewWriter(ctx)

	if _, err = io.Copy(wc, r); err != nil {
		return nil, nil, err
	}

	if err = wc.Close(); err != nil {
		return nil, nil, err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)
	log.Printf("Post is saved to GCS: %s\n", attrs.MediaLink)
	return obj, attrs, err
}