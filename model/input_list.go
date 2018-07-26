
package model

import "gopkg.in/mgo.v2/bson"

type InList struct {
	TaskID          	bson.ObjectId `bson:"task_id"       json:"task_id"`
	Cate				string        `bson:"cate"       	json:"cate"`
	Each				[]InEach	  `bson:"each"       	json:"each"`
}
