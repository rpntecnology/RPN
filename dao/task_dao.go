package dao

import (
	"gopkg.in/mgo.v2"
	"log"
	"RPN/config"
	"fmt"
	"RPN/model"
	"gopkg.in/mgo.v2/bson"
)

type TaskDAO struct {
	Server   string
	Database string
	db       *mgo.Database
	session  *mgo.Session
	err      error
}

const (
	TASK_COLLECTION = "task"
)


// connect to database
func (m *TaskDAO) Connect() {
	m.Server = config.DB_SERVER
	m.Database = config.DB_NAME
	log.Println(m.Server)
	log.Println(m.Database)
	m.session,m.err = mgo.Dial(m.Server)
	if m.err != nil {
		log.Fatal(m.err)
	}
	m.db = m.session.DB(m.Database)
}

func (m *TaskDAO) AddTask(task model.Task) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Insert(&task)
	fmt.Printf("%+v\n", task)
	return err
}

func (m *TaskDAO) AddPrevImage(imageSlot model.ImageSlot) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Update(
		bson.M{"task_id": imageSlot.TaskID, "item_list": bson.M{"$elemMatch": bson.M{"cate": imageSlot.Cate, "item":imageSlot.ItemId}}},
		bson.M{"$push": bson.M{"item_list.$.before":imageSlot}})


	m.db.C(TASK_COLLECTION).Update(
		bson.M{"task_id": imageSlot.TaskID},
		bson.M{"$inc": bson.M{"totalImage": 1}})
	return err
}

//func (m* TaskDAO) FindTotalImage(taskID bson.ObjectId) error {
//	m.Connect()
//	defer m.session.Close()
//	m.db.C(TASK_COLLECTION).Update(
//		bson.M{"task_id": taskID},
//		bson.M{""}
//	)
//}

func (m *TaskDAO) FindTasksByUsername(username string) (error, []model.Task) {
	m.Connect()
	defer m.session.Close()
	var tasks []model.Task
	err := m.db.C(TASK_COLLECTION).Find(bson.M{"username": username}).All(&tasks)
	return err, tasks
}

func (m *TaskDAO) FindById(taskId bson.ObjectId) (error, model.Task) {
	m.Connect()
	defer m.session.Close()
	var task model.Task
	err := m.db.C(TASK_COLLECTION).Find(bson.M{"task_id": taskId}).One(&task)
	return err, task
}


//func (m *TaskDAO) FindImageByCategory(taskID bson.ObjectId, category string) (error, []model.ImageSlot) {
//	m.Connect()
//	defer m.session.Close()
//	var Images []model.ImageSlot
//	pipeline := []bson.M {
//		bson.M{"$match": bson.M{"task_id": taskID}},
//		bson.M{"$unwind": "$image"},
//		bson.M{"$match": bson.M{"image.category": category}},
//	}
//	pipe := m.db.C(TASK_COLLECTION).Pipe(pipeline)
//	log.Println(pipeline)
//	resp := []bson.M{}
//	var task1 model.Task
//	//resp := []model.Task{}
//	err1 := pipe.All(&resp)
//	err2 := pipe.One(&task1)
//	if err1 != nil {
//		log.Println(err1.Error())
//	}
//	if err2 != nil {
//		log.Println(err1.Error())
//	}
//	log.Println("resp: ")
//	log.Println(resp[1])
//	log.Println(task1)
//	err := pipe.All(&Images)
//	return err, Images
//}

//new added functions
func (m *TaskDAO) DeleteImageByImageID(taskID bson.ObjectId, imageID bson.ObjectId) (error) {
	m.Connect()
	defer m.session.Close()
	//var task model.Task
	err := m.db.C(TASK_COLLECTION).Update(bson.M{"task_id": taskID}, bson.M{"$pull": bson.M{"image" : bson.M{"image_id" : imageID}}})
	return err
}

func (m *TaskDAO) DeleteTaskByTaskID(taskID bson.ObjectId) (error) {
	m.Connect()
	defer m.session.Close()
	//db.collection.remove({_id: item._id})
	err := m.db.C(TASK_COLLECTION).Remove(bson.M{"task_id" : taskID})
	return err
}

func (m *TaskDAO) AssignTaskToAnotherUser(taskID bson.ObjectId, newUser string) error {
	m.Connect()
	defer m.session.Close()
	//var task model.Task
	err := m.db.C(TASK_COLLECTION).Update(bson.M{"task_id": taskID}, bson.M{"$set": bson.M{"username" : newUser}})
	return err
}

func (m *TaskDAO) AddCategory(taskID bson.ObjectId, cate model.List) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Update(bson.M{"task_id": taskID}, bson.M{"$push": bson.M{"list": cate}})
	return err
}

func (m *TaskDAO) AddItem(taskID bson.ObjectId, cate string, item model.Each) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Update(
		bson.M{"task_id": taskID, "list.cate": cate},
		bson.M{"$push": bson.M{"list.$.each": item}})
	return err
}

func (m *TaskDAO) FindImgs(taskID bson.ObjectId, cate, item, status string) (error, model.Task) {
	m.Connect()
	defer m.session.Close()
	var task model.Task
	//err := m.db.C(TASK_COLLECTION).Update(
	//	bson.M{"task_id": imageSlot.TaskID, "list.cate": imageSlot.Cate, "list.each.item": imageSlot.ItemId},
	//	bson.M{"$push": bson.M{"list.0.each.$.before":imageSlot}})
	err := m.db.C(TASK_COLLECTION).Find(
		bson.M{"task_id": taskID, "list.cate": cate, "list.0.each.item": item}).One(&task)

	return err, task
}

