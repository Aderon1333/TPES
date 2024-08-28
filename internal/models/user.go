package models

type User struct {
	// gorm.Model
	Login    string `json:"login" gorm:"unique"`
	Password string `json:"password"`
}
