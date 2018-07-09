package dao

import (
	"gopkg.in/mgo.v2"
	"log"
	"RPN/model"
	"RPN/config"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"crypto/tls"
	"net"
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

	dialInfo, err := mgo.ParseURL(m.Server)
	if err != nil {
		log.Print("error in parsing url")
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)

		if err != nil {
			log.Print("Error in tls dial")
			log.Print(err.Error())
		}
		return conn, err
	}

	m.session, m.err = mgo.DialWithInfo(dialInfo)

	//m.session,m.err = mgo.Dial(m.Server)
	if m.err != nil {
		log.Print("error in connceting to mongodb")
		log.Fatal(m.err)
	}
	m.db = m.session.DB(m.Database)
	log.Print("successfully connect to db")
}

//Add User
func (m *UserDAO) AddUser(user model.User) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(USER_COLLECTION).Insert(&user)
	fmt.Printf("%+v\n", user)
	return err
}

func (m *UserDAO) CheckUser(username, password string) bool{
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
	//if err != nil {
	//	log.Println("Error in finding username: " + username)
	//	return err, user
	//}
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

func (m *UserDAO) DeleteUser(username string) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(USER_COLLECTION).Remove(bson.M{"username": username})
	return err
}

func (m *UserDAO) AssignTask(username string, taskId bson.ObjectId) error {
	m.Connect()
	defer m.session.Close()
	err := m.db.C(USER_COLLECTION).Update(
		bson.M{"username": username},
		bson.M{"$push": bson.M{"task_ids": taskId}})
	return err
}