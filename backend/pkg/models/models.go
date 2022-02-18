package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Image    string `json:"img"`
	UserName string `gorm:"column:username" json:"username"`
	Email    string `gorm:"column:email" json:"email"`
	Pass     string `gorm:"column:password" json:"password"`
	Desc     string `gorm:"column:desc" json:"desc"`
}
