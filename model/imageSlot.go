package model

import (
	"gopkg.in/mgo.v2/bson"
)

type ImageSlot struct {
	ImageID			bson.ObjectId `bson:"image_id" json:"image_id"`
	Category		string        `bson:"category" json:"category"`
	URL				string        `bson:"url"      json:"url"`
	Status			string        `bson:"status" json:"status"`
}
