package model

import "gopkg.in/mgo.v2/bson"

type Image struct {
	ID          		bson.ObjectId `bson:"_id"       	json:"id"`
	Name    			string        `bson:"name"      	json:"name"`
	Src					string		  `bson:"src"      		json:"src"`
}
