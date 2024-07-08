package core

type Task struct {
	ID     int64  `json:"id" bson:"_id"`
	Status string `json:"status" bson:"status"`
	Item   string `json:"item" bson:"item"`
}
