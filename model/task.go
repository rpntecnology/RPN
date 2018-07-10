package model

import "gopkg.in/mgo.v2/bson"

type Task struct {
	TaskID          bson.ObjectId `bson:"task_id" 			json:"task_id"`
	AssetNumber     string        `bson:"asset_num"         json:"asset_num"`
	Invoice         string        `bson:"invoice" 			json:"invoice"`
	BillTo			string		  `bson:"billTo" 			json:"billTo"`
	StartDate       string        `bson:"startDate" 		json:"startDate"`
	CompletionDate  string	  	  `bson:"completionDate" 	json:"completionDate"`
	InvoiceDate		string        `bson:"invoiceDate" 	    json:"invoiceDate"`
	Username        string        `bson:"username" 			json:"username"`
	Name			string        `bson:"name" 				json:"name"`
	Address         string        `bson:"address" 			json:"address"`
	City			string        `bson:"city" 			    json:"city"`
	Year			string		  `bson:"year" 				json:"year"`
	Stories			string        `bson:"stories" 			json:"stories"`
	Area            string        `bson:"area" 			    json:"area"`
	TotalCost       string        `bson:"totolCost" 	    json:"totalCost"`
	ItemList		[]Each        `bson:"item_list" 		json:"item_list"`
	TotalImage      int           `bson:"totalImage" 		json:"totalImage"`
	Stage           string        `bson:"stage" 			json:"stage"`
}