package dao

import (
	"gopkg.in/mgo.v2"
	"log"
	"RPN/model"
	"RPN/config"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type UserDAO struct {
	Server   string
	Database string
}
var db *mgo.Database

const (
	COLLECTION = "user"
)

// connect to database
func (m *UserDAO) Connect() {
	m.Server = config.DB_SERVER
	m.Database = config.DB_NAME
	log.Println(m.Server)
	log.Println(m.Database)
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

//Add User
func (m *UserDAO) AddUser(user model.User) error {
	m.Connect()
	err := db.C(COLLECTION).Insert(&user)
	fmt.Printf("%+v\n", user)
	return err
}

func (m *UserDAO) FindAll() ([]model.User, error) {
	m.Connect()
	var users []model.User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}