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
	db       *mgo.Database
	session  *mgo.Session
	err      error
}

const (
	USER_COLLECTION = "user"
)



// connect to database
func (m *UserDAO) Connect() {
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

//Add User
func (m *UserDAO) AddUser(user model.User) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(USER_COLLECTION).Insert(&user)
	fmt.Printf("%+v\n", user)
	return err
}

func (m *UserDAO) CheckUser(username, password string) bool {
	m.Connect()
	defer m.session.Close()
	var user model.User
	err := m.db.C(USER_COLLECTION).Find(bson.M{"username": username}).One(&user)
	if err != nil {
		log.Println("Error in finding username: " + username)
		return false
	}

	return username == user.Username && password == user.Password
}

func (m *UserDAO) FindUser(username string) (error, model.User) {
	m.Connect()
	defer m.session.Close()
	var user model.User
	err := m.db.C(USER_COLLECTION).Find(bson.M{"username": username}).One(&user)
	if err != nil {
		log.Println("Error in finding username: " + username)
		return err, user
	}
	return err, user
}


func (m *UserDAO) FindAll() ([]model.User, error) {
	m.Connect()
	defer m.session.Close()
	var users []model.User
	err := m.db.C(USER_COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

func (m *UserDAO) FindAuthority(username string) string {
	m.Connect()
	defer m.session.Close()
	var user model.User
	err := m.db.C(USER_COLLECTION).Find(bson.M{"username": username}).One(&user)
	if err != nil {
		log.Println("Error in finding username: " + username)
		return ""
	}
	return user.Authority
}