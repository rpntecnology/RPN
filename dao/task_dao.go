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

func (m *TaskDAO) AddImage(taskId bson.ObjectId, imageSlot model.ImageSlot) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Update(
		bson.M{"task_id": taskId},
		bson.M{"$push": bson.M{"image": imageSlot}})
	return err
}

func (m *TaskDAO) FindByUsername() ([]model.Task, error) {
	m.Connect()
	var tasks []model.Task
	err := m.db.C(TASK_COLLECTION).Find(bson.M{}).All(&tasks)
	return tasks, err
}

func (m *TaskDAO) FindImageByCategory(taskID bson.ObjectId, Category string) (error, []bson.ObjectId) {
	m.Connect()
	var imageIDs []bson.ObjectId
	//var task model.Task
	err := m.db.C(TASK_COLLECTION).Find(bson.M{"TaskID": taskID}).Select(bson.M{"Image": bson.M{"$elemMatch": bson.M{"Category": Category}}}).One(&imageIDs)
	if err != nil {
		log.Println("Error in finding task with this taskID: " + taskID)
		return err, nil
	}

	return err, imageIDs
}
