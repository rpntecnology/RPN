package model

type InEach struct {
	//TaskID          	bson.ObjectId `bson:"task_id"       json:"task_id"`
	Cate				string		  `bson:"cate"  		json:"cate"`
	Item    			int        	  `bson:"item"      	json:"item"`
	Description         string        `bson:"description"   json:"description"`
	Qty        			int           `bson:"qty"  			json:"qty"`
	UM					string		  `bson:"UM"  			json:"UM"`
	PPU 				float64		  `bson:"PPU"  			json:"PPU"`
	Cost				float64		  `bson:"cost"  		json:"cost"`
	Amount              float64        `bson:"amount"  		json:"amount"`
	Taxable             bool		  `bson:"taxable"  		json:"taxble"`
	Tax                 float64		  `bson:"tax"  			json:"tax"`
	Comments 			string        `bson:"comments"  	json:"comments"`
	Image				[]ImageSlot	  `bson:"image"  		json:"image"`
}
