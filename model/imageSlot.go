package model

import (
	"gopkg.in/mgo.v2/bson"
)

type ImageSlot struct {
	ImageID			bson.ObjectId `bson:"image_id" 		json:"image_id"`
	TaskID			bson.ObjectId `bson:"task_id" 		json:"task_id"`
	Name			string        `bson:"name"  		json:"name"`
	Format			string		  `bson:"format" 		json:"format"`
	Cate		    string        `bson:"cate" 			json:"cate"`
	ItemId			string		  `bson:"itemId         json:"itemId""`
	Src				string        `bson:"src"      		json:"src"`
	Status			string        `bson:"status" 		json:"status"`
}
