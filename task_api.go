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

func InitTaskHandler(w http.ResponseWriter, r *http.Request) {
	if CheckAuth(r) < AUTH_TO_MANAGE_TASK {
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
	task.Stage = "0"
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

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if CheckAuth(r) < AUTH_TO_MANAGE_TASK {
		log.Println("No authority to add task")
		respondWithError(w, http.StatusForbidden, "No authority to add task")
		return
	}

	log.Println("Received new add task request")
	defer r.Body.Close()
	var inTask model.InputTask
	if err := json.NewDecoder(r.Body).Decode(&inTask); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	task := TransformTask(inTask)
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

func ParseJsonHandler(w http.ResponseWriter, r *http.Request) {
	if CheckAuth(r) < AUTH_TO_MANAGE_TASK {
		log.Println("No authority to add task")
		respondWithError(w, http.StatusForbidden, "No authority to add task")
		return
	}

	log.Println("Received new add task request")
	defer r.Body.Close()
	var inTask model.InputTask
	if err := json.NewDecoder(r.Body).Decode(&inTask); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//log.Println(inTask)
	//log.Println("intask id: " + inTask.TaskID)
	//log.Println("username: " + inTask.Username)
	task := TransformTask(inTask)
	//task.TaskID = bson.NewObjectId()

	if err, existTask := taskDao.FindById(task.TaskID); err != nil {
		log.Println("DB find error: no such task")
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		task.AssetNumber = existTask.AssetNumber
		task.StartDate = existTask.StartDate
		task.CompletionDate = existTask.CompletionDate
		//task.Stage = existTask.Stage
	}

	if err := taskDao.UpdateTask(task); err != nil {
		log.Println("DB insert error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondWithJson(w, http.StatusCreated, task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if CheckAuth(r) < AUTH_TO_MANAGE_TASK {
		log.Println("No authority to update task")
		respondWithError(w, http.StatusForbidden, "No authority to update task")
		return
	}

	log.Println("Received new update task request")
	defer r.Body.Close()
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//function _update($id, $data, $options=array()){
	//	$temp = array();
	//	foreach($data as $key => $value) {
	//		$temp["some_key.".$key] = $value;
	//		}
	//		$collection->update( array('_id' => $id), array('$set' => $temp) );
	//		}
	//
	//		_update('1', array('param2' => 'some data'));

	taskId := bson.ObjectIdHex(r.URL.Query().Get("task_id"))
	task.TaskID = taskId
	//og.Print(task.TaskID)
	if err := taskDao.UpdateTask(task); err != nil {
		log.Println("DB insert error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondWithJson(w, http.StatusCreated, task)
}

func FindTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	log.Println("Received find one task request")

	defer r.Body.Close()
	taskId := bson.ObjectIdHex(r.URL.Query().Get("task_id"))
	err, tasks := taskDao.FindById(taskId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tasks)
}

func AddImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	log.Println("Received new add image request")
	if CheckAuth(r) < 0 {
		log.Println("No authority to post images")
		respondWithError(w, http.StatusForbidden, "No authority to post images")
		return
	}

	defer r.Body.Close()
	var imageSlot model.ImageSlot

	// 32 << 20 is the maxMemory param for ParseMultipartForm, equals to 32MB (1MB = 1024 * 1024 bytes = 2^20 bytes)
	// After you call ParseMultipartForm, the file will be saved in the server memory with maxMemory size.
	// If the file size is larger than maxMemory, the rest of the data will be saved in a system temporary file.
	r.ParseMultipartForm(32 << 20)
	taskId := r.FormValue("task_id")
	name := r.FormValue("name")
	cate := r.FormValue("cate")
	itemId,_  := strconv.Atoi(r.FormValue("item_id"))
	status := r.FormValue("status")

	log.Println("itemId: " + strconv.Itoa(itemId))
	log.Println("taskId: " + taskId)
	log.Println("cate: " + cate)
	log.Println("status" + status)

	imageSlot.ImageID = bson.NewObjectId()
	imageSlot.TaskID = bson.ObjectIdHex(taskId)
	imageSlot.Name = name
	imageSlot.Cate = cate
	imageSlot.ItemId = itemId
	imageSlot.Status = status

	file, _, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Image is not available")
		log.Println("Image is not available, err: " + err.Error())
		return
	}
	defer file.Close()

	// parse multiple images
	//fhs := r.MultipartForm.File["images"]
	//for _, fh := range fhs {
	//
	//}
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

	log.Println(imageSlot)

	if err := taskDao.AddImage(imageSlot); err != nil {
		respondWithError(w, http.StatusConflict, err.Error())
		log.Println("DB insert error, err: " + err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, imageSlot)
}

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Received new upload image request")

	defer r.Body.Close()
	var imageSlot model.ImageSlot

	// 32 << 20 is the maxMemory param for ParseMultipartForm, equals to 32MB (1MB = 1024 * 1024 bytes = 2^20 bytes)
	// After you call ParseMultipartForm, the file will be saved in the server memory with maxMemory size.
	// If the file size is larger than maxMemory, the rest of the data will be saved in a system temporary file.
	r.ParseMultipartForm(32 << 20)
	imageSlot.ImageID = bson.NewObjectId()

	file, _, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Image is not available")
		log.Println("Image is not available, err: " + err.Error())
		return
	}
	defer file.Close()

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
	log.Println(imageSlot)
	respondWithJson(w, http.StatusCreated, imageSlot.Src)
}

func FindAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	log.Println("Received find all task request")
	if CheckAuth(r) < AUTH_TO_DELETE {
		log.Println("No authority to view all tasks")
		respondWithError(w, http.StatusForbidden, "No authority to view all tasks")
		return
	}

	tasks, err := taskDao.FindAllTasks()

 	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tasks)
}


func FindImgURLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	log.Println("Received Find images urls request")

	taskId := r.URL.Query().Get("task_id")
	cate := r.URL.Query().Get("cate")
	itemId, _ := strconv.Atoi(r.URL.Query().Get("item_id"))
	status := r.URL.Query().Get("status")

	log.Println("task_id: " + taskId)
	log.Println("cate: " + cate)
	log.Println("status: " + status)

	err, task := taskDao.FindById(bson.ObjectIdHex(taskId))
	if err != nil {
		log.Println("Error in finding task")
		respondWithError(w, http.StatusInternalServerError, "Error in finding task")
		return
	}

	var imgs []model.ImageSlot
	for _, item := range task.ItemList {
		if item.Cate == cate && item.Item == itemId {
			switch status {
			case "before":
				imgs = item.Before
			case "during":
				imgs = item.During
			case "after":
				imgs = item.After
			}
		}
	}

	if len(imgs) == 0 {
		log.Println("Error in finding images in such item: " + strconv.Itoa(itemId))
		respondWithError(w, http.StatusInternalServerError, "Error in finding images in such item: " + strconv.Itoa(itemId))
		return
	}

	log.Println(imgs)
	respondWithJson(w, http.StatusOK, imgs)
}

func FindOneItemImgURLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	log.Println("Received Find images urls of one item request")

	taskId := r.URL.Query().Get("task_id")
	cate := r.URL.Query().Get("cate")
	itemId, _ := strconv.Atoi(r.URL.Query().Get("item_id"))
	log.Println("task_id: " + taskId)
	log.Println("cate: " + cate)

	err, task := taskDao.FindById(bson.ObjectIdHex(taskId))
	if err != nil {
		log.Println("Error in finding task")
		respondWithError(w, http.StatusInternalServerError, "Error in finding task")
		return
	}

	var imgs map[string][]model.ImageSlot
	imgs = make(map[string][]model.ImageSlot)
	for _, item := range task.ItemList {
		if item.Cate == cate && item.Item == itemId {
			imgs["before"] = item.Before
			imgs["during"] = item.During
			imgs["after"] = item.After
		}
	}

	if len(imgs) == 0 {
		log.Println("Error in finding images in such item: " + strconv.Itoa(itemId))
		respondWithError(w, http.StatusInternalServerError, "Error in finding images in such item: " + strconv.Itoa(itemId))
		return
	}

	log.Println(imgs)
	respondWithJson(w, http.StatusOK, imgs)
}

func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	if CheckAuth(r) < 0 {
		respondWithError(w, http.StatusInternalServerError, "No authority to delete image")
		log.Println("No authority to delete image")
		return
	}
	log.Println("Received delete image request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	taskId := r.URL.Query().Get("task_id")
	imageId := r.URL.Query().Get("image_id")
	log.Println(taskId)
	log.Println(imageId)

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
	if CheckAuth(r) < AUTH_TO_DELETE {
		respondWithError(w, http.StatusInternalServerError, "No authority to delete task")
		log.Println("No authority to delete task")
		return
	}
	log.Println("Received delete task request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	taskId := r.URL.Query().Get("task_id")
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

func FinishTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	log.Println("Received finish a task request")

	if CheckAuth(r) < AUTH_TO_MANAGE_TASK {
		respondWithError(w, http.StatusInternalServerError, "No authority to finish a task")
		log.Println("No authority to finish a task")
		return
	}

	taskId := r.URL.Query().Get("task_id")
	err := taskDao.FinishTask(bson.ObjectIdHex(taskId))

	if err != nil {
		log.Println("taskID: " + taskId + " DB find error")
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, "success")
}

func TerminateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	log.Println("Received terminate a task request")

	if CheckAuth(r) < AUTH_TO_DELETE {
		respondWithError(w, http.StatusInternalServerError, "No authority to terminate a task")
		log.Println("No authority to terminate a task")
		return
	}

	taskId := r.FormValue("task_id")
	errorStage := r.FormValue("error_stage")
	reason := r.FormValue("reason")

	err := taskDao.TerminateTask(bson.ObjectIdHex(taskId), errorStage, reason)
	if err != nil {
		log.Print("taskId: " + taskId)
		log.Print("error stage: " + errorStage)
		log.Print("reason: " + reason)
		log.Println(err.Error())
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, "success")
}

//func ChangeContractorHandler(w http.ResponseWriter, r *http.Request) {
//	if CheckAuth(r) < AUTH_TO_MANAGE_TASK {
//		respondWithError(w, http.StatusInternalServerError, "No authority to assign work")
//		log.Println("No authority to assign work")
//		return
//	}
//	log.Println("Received change task's contractor request")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Content-Type", "application/json")
//	taskId, _ := r.URL.Query().Get("task_id"), 64
//	log.Println(taskId)
//	user, _ := r.URL.Query().Get("username"), 64
//	log.Println(user)
//
//	if err := taskDao.AssignTaskToAnotherUser(bson.ObjectIdHex(taskId), user); err != nil {
//		log.Println("taskID: "+taskId +" user: " + user)
//		log.Println("DB find error")
//		log.Println(err.Error())
//		respondWithError(w, http.StatusConflict, err.Error())
//		return
//	}
//	log.Println("changed task's contractor successfully")
//}

func AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	//if checkAuth(r) < AUTH_TO_MANAGE_TASK {
	//	respondWithError(w, http.StatusInternalServerError, "No authority to add category")
	//	log.Println("No authority to add category")
	//	return
	//}
	//log.Println("Received add category request")
	//defer r.Body.Close()
	//var cate model.List
	//if err := json.NewDecoder(r.Body).Decode(&cate); err != nil {
	//	respondWithError(w, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//
	//if err := taskDao.AddCategory(cate.TaskID, cate); err != nil {
	//	log.Println("DB insert error")
	//	log.Println(err.Error())
	//	respondWithError(w, http.StatusConflict, err.Error())
	//	return
	//}
	//w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//respondWithJson(w, http.StatusCreated, cate)
}

//func AddItemHandler (w http.ResponseWriter, r *http.Request) {
//	if CheckAuth(r) < AUTH_TO_MANAGE_TASK {
//		respondWithError(w, http.StatusInternalServerError, "No authority to add items")
//		log.Println("No authority to add items")
//		return
//	}
//	log.Println("Received add items request")
//	defer r.Body.Close()
//	var item model.Each
//	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
//		respondWithError(w, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	if err := taskDao.AddItem(item.TaskID, item.Cate, item); err != nil {
//		log.Println("DB insert error")
//		log.Println(err.Error())
//		respondWithError(w, http.StatusConflict, err.Error())
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	respondWithJson(w, http.StatusCreated, item)
//}

func CheckAuth(r *http.Request) int {
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	authority := claims.(jwt.MapClaims)["authority"]
	auth, _ := strconv.Atoi(authority.(string))
	log.Print("auth: ")
	log.Println(auth)
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

func TransformTask(inTask model.InputTask) model.Task {
	var task model.Task
	task.TaskID = bson.ObjectIdHex(inTask.TaskID)
	task.Invoice = inTask.Invoice
	task.BillTo = inTask.BillTo
	task.CompletionDate = inTask.CompletionDate
	task.InvoiceDate = inTask.InvoiceDate
	task.Username = inTask.Username
	task.Name = inTask.Name
	task.Address = inTask.Address
	task.City = inTask.City
	task.Year = inTask.Year
	task.Stories = inTask.Stories
	task.Area = inTask.Area
	task.TotalCost = inTask.TotalCost
	task.TotalImage = inTask.TotalImage
	task.Stage = "1"
	task.ErrorStage = ""
	task.ErrorReason = ""

	var itemList []model.Each

	cateList := inTask.List
	for _, cate := range cateList {
		for _, item := range cate.Each {
			var myItem model.Each
			myItem.Cate = cate.Cate
			for _, img := range item.Image {
				var bImg model.ImageSlot
				bImg.ImageID = bson.NewObjectId()
				bImg.TaskID = bson.ObjectIdHex(inTask.TaskID)
				bImg.Src = img.Src
				bImg.Name = img.Name
				myItem.Before = append(myItem.Before, bImg)
			}
			itemList = append(itemList, myItem)
		}
	}
	task.ItemList = itemList
	return task
}