package model

import "gopkg.in/mgo.v2/bson"

type ResponseProfile struct {
	TaskID          bson.ObjectId `bson:"task_id" 			json:"task_id"`
	Username        string        `bson:"username" 			json:"username"`
	Name			string        `bson:"name" 				json:"name"`
	Invoice         string        `bson:"invoice" 			json:"invoice"`
	BillTo			string		  `bson:"billTo" 			json:"billTo"`
	CompletionDate  string	  	  `bson:"completionDate" 	json:"completionDate"`
	InvoiceDate		string        `bson:"invoiceDate" 	    json:"invoiceDate"`
	Address         string        `bson:"address" 			json:"address"`
	City			string        `bson:"city" 			    json:"city"`
	Year			string		  `bson:"year" 				json:"year"`
	Stories			string        `bson:"stories" 			json:"stories"`
	Area            string        `bson:"area" 			    json:"area"`
	TotalCost       string        `bson:"totolCost" 	    json:"totalCost"`
	ItemList       	[]Each        `bson:"list" 		        json:"list"`
	TotalImage      int           `bson:"totalImage" 		json:"totalImage"`
	Stage           string        `bson:"stage" 			json:"stage"`
}
