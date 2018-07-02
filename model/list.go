package model

import "gopkg.in/mgo.v2/bson"

type List struct {
	TaskID          	bson.ObjectId `bson:"task_id"       json:"task_id"`
	Cate				string        `bson:"cate"       	json:"cate"`
	Each				[]Each		  `bson:"each"       	json:"each"`
}