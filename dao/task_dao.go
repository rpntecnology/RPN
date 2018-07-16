package dao

import (
	"gopkg.in/mgo.v2"
	"log"
	"RPN/config"
	"fmt"
	"RPN/model"
	"gopkg.in/mgo.v2/bson"
	"crypto/tls"
	"net"
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

	dialInfo, err := mgo.ParseURL(m.Server)
	if err != nil {
		log.Print("error in parsing url")
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	m.session, m.err = mgo.DialWithInfo(dialInfo)

	//m.session,m.err = mgo.Dial(m.Server)
	if m.err != nil {
		log.Print("error in connceting to mongodb")
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

func (m *TaskDAO) UpdateTask(task model.Task) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Update(
		bson.M{"task_id": task.TaskID}, &task)
	return err
}

func (m *TaskDAO) AddPrevImage(imageSlot model.ImageSlot) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Update(
		bson.M{"task_id": imageSlot.TaskID, "item_list": bson.M{"$elemMatch": bson.M{"cate": imageSlot.Cate, "item":imageSlot.ItemId}}},
		bson.M{"$push": bson.M{"item_list.$.before":imageSlot}})

	// update total images
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

func (m *TaskDAO) FindByIdPreview(taskId bson.ObjectId) (error, model.Task) {
	m.Connect()
	defer m.session.Close()
	var task model.Task
	err := m.db.C(TASK_COLLECTION).Find(bson.M{"task_id": taskId}).
		Select(bson.M{"task_id":1, "asset_num":1, "startDate":1,
		"completionDate":1, "username":1, "name":1, "address":1,
		"city":1}).One(&task)
	return err, task
}

func (m *TaskDAO) FindAllTasks() ([]model.Task, error) {
	m.Connect()
	defer m.session.Close()
	var tasks []model.Task
	//err := m.db.C(USER_COLLECTION).Find(bson.M{}).Select(bson.M{"_id":0, "password":0}).All(&users)
	err := m.db.C(TASK_COLLECTION).Find(bson.M{}).
		Select(bson.M{"task_id":1, "asset_num":1, "startDate":1,
		"completionDate":1, "username":1, "name":1, "address":1,
		"city":1}).All(&tasks)
	return tasks, err
}

//new added functions
func (m *TaskDAO) DeleteImageByImageID(taskID bson.ObjectId, imageID bson.ObjectId) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Update(bson.M{"task_id": taskID}, bson.M{"$pull": bson.M{"image" : bson.M{"image_id" : imageID}}})

	// update total images
	m.db.C(TASK_COLLECTION).Update(
		bson.M{"task_id": taskID},
		bson.M{"$inc": bson.M{"totalImage": -1}})
	return err
	return err
}

func (m *TaskDAO) DeleteTaskByTaskID(taskID bson.ObjectId) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Remove(bson.M{"task_id" : taskID})
	return err
}

func (m *TaskDAO) AssignUserToTask(taskID bson.ObjectId, newUser string) error {
	m.Connect()
	defer m.session.Close()
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



