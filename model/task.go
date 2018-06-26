package model

import "gopkg.in/mgo.v2/bson"

type Task struct {
	TaskID          bson.ObjectId `bson:"task_id" json:"task_id"`
	TaskName        string        `bson:"taskname" json:"taskname"`
	Username        string        `bson:"username" json:"username"`
	Address         string        `bson:"address" json:"address"`
	Category        []string      `bson:"category" json:"category"`
	Image           []ImageSlot   `bson:"image" json:"image"`
	IsDone          string        `bson:"isDone" json:"isDone"`
}