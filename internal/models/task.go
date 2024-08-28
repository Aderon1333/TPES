package models

type Task struct {
	//gorm.Model
	ID     int64  `json:"id" bson:"_id" gorm:"unique"`
	Status string `json:"status" bson:"status"`
	Item   string `json:"item" bson:"item"`
}
