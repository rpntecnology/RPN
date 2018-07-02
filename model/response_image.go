package model

type ResponseImage struct {
	Name    			string        `bson:"name"      	json:"name"`
	Src					string		  `bson:"src"      		json:"src"`
}
