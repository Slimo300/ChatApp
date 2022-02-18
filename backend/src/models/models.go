package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Image    string    `json:"img"`
	UserName string    `gorm:"column:username" json:"username"`
	Email    string    `gorm:"column:email" json:"email"`
	Pass     string    `gorm:"column:password" json:"password"`
	Desc     string    `gorm:"column:desc" json:"desc"`
	Active   time.Time `gorm:"column:activity" json:"activity"`
	SignUp   time.Time `gorm:"column:signup" json:"signup"`
	LoggedIn bool      `gorm:"column:logged" json:"logged"`
}

type Group struct {
	gorm.Model
	Image   string    `json:"img"`
	Name    string    `gorm:"column:name" json:"name"`
	Created time.Time `gorm:"column:created" json:"created"`
}

type Message struct {
	gorm.Model
	Posted time.Time `gorm:"column:posted" json:"posted"`
}

type Priv struct {
	Add_user    bool
	Delete_user bool
	Settings    bool
}
