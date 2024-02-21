package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Age      uint32 `json:"age"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}
