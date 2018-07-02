package model

import "gopkg.in/mgo.v2/bson"

type Each struct {
	TaskID          	bson.ObjectId `bson:"task_id"       json:"task_id"`
	Cate				string		  `bson:"cate"  		json:"cate"`
	Item    			string        `bson:"item"      	json:"item"`
	Description         string        `bson:"description"   json:"description"`
	Qty        			string        `bson:"qty"  			json:"qty"`
	UM					string		  `bson:"UM"  			json:"UM"`
	PPU 				string		  `bson:"PPU"  			json:"PPU"`
	Cost				string		  `bson:"cost"  		json:"cost"`
	Amount              string        `bson:"amount"  		json:"amount"`
	Taxable             bool		  `bson:"taxable"  		json:"taxble"`
	Tax                 string		  `bson:"tax"  			json:"tax"`
	Comments 			string        `bson:"comments"  	json:"comments"`
	Before				[]ImageSlot	  `bson:"before"  		json:"before"`
	During				[]ImageSlot	  `bson:"during"  		json:"during"`
	After 				[]ImageSlot   `bson:"after"  		json:"after"`
}
