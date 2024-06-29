package core

type Task struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
	Item   string `json:"item"`
}
