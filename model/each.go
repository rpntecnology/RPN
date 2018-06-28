package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Each struct {
	ID          		bson.ObjectId `bson:"each_id"       	json:"each_id"`
	Item    			int        	  `bson:"item"      	json:"item"`
	Description         string        `bson:"description"   json:"description"`
	Qty        			int       	  `bson:"qty"  			json:"qty"`
	UM					string		  `bson:"UM"  			json:"UM"`
	PPU 				float64		  `bson:"PPU"  			json:"PPU"`
	Cost				float64		  `bson:"cost"  		json:"cost"`
	Amount              float64       `bson:"amount"  		json:"amount"`
	Taxable             bool		  `bson:"amount"  		json:"amount"`
	Tax                 float64		  `bson:"amount"  		json:"amount"`
	Comments 			string        `bson:"comments"  	json:"comments"`
	Before				[]ImageSlot	  `bson:"before"  		json:"before"`
	During				[]ImageSlot	  `bson:"during"  		json:"during"`
	After 				[]ImageSlot       `bson:"after"  		json:"after"`
}
