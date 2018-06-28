package model

type List struct {
	//ID          		bson.ObjectId `bson:"list_id"       json:"list_id"`
	Cate				string        `bson:"cate"       	json:"cate"`
	Each				[]Each		  `bson:"each"       	json:"each"`
}