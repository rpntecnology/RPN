package dao

import (
	"gopkg.in/mgo.v2"
	"log"
	"RPN/config"
)

type TaskDAO struct {
	Server   string
	Database string
	db       *mgo.Database
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
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	m.db = session.DB(m.Database)
}

func (m *TaskDAO) AddTask() {

}