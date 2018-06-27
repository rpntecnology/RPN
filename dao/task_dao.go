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

func (m *TaskDAO) AddImage(imageSlot model.ImageSlot) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(TASK_COLLECTION).Update(
		bson.M{"task_id": imageSlot.TaskID},
		bson.M{"$push": bson.M{"image": imageSlot}})
	return err
}

func (m *TaskDAO) FindByUsername() ([]model.Task, error) {
	m.Connect()
	var tasks []model.Task
	err := m.db.C(TASK_COLLECTION).Find(bson.M{}).All(&tasks)
	return tasks, err
}

func (m *TaskDAO) FindImageByCategory(taskID bson.ObjectId, Category string) (error, []model.ImageSlot) {
	m.Connect()
	var Images []model.ImageSlot
	//var task model.Task
	err := m.db.C(TASK_COLLECTION).Find(bson.M{"task_id": taskID}).Select(bson.M{"image": bson.M{"$elemMatch": bson.M{"category": Category}}}).All(&Images)
	if err != nil {
		log.Println("Error in finding task with this taskID: " + taskID)
		return err, Images
	}

	return err, Images
}

func (m *TaskDAO) FindImageByCategoryII(taskID bson.ObjectId, category string) (error, []model.ImageSlot) {
	m.Connect()
	var Images []model.ImageSlot
	pipeline := []bson.M {
		bson.M{"$match": bson.M{"task_id": taskID}},
		bson.M{"$unwind": "$image"},
		bson.M{"$match": bson.M{"image.category": category}},
	}
	pipe := m.db.C(TASK_COLLECTION).Pipe(pipeline)
	resp := []bson.M{}
	err1 := pipe.All(&resp)
	if err1 != nil {
		log.Println(err1.Error())
	}
	log.Println(resp)

	err := pipe.All(&Images)
	return err, Images
}
