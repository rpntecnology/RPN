package model

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID          bson.ObjectId `bson:"_id"       json:"id"`
	Username    string        `bson:"username"  json:"username"`
	Password    string        `bson:"password"  json:"password"`
	Firstname   string        `bson:"firstname" json:"firstname"`
	Lastname    string        `bson:"lastname"  json:"lastname"`
	Email	    string        `bson:"email"     json:"email"`
	Phone       string        `bson:"phone"     json:"phone"`
	Authority   string        `bson:"authority" json:"authority"`
}