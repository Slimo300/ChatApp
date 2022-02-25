package models

import (
	"time"
)

type User struct {
	ID       int       `gorm:"primaryKey"`
	UserName string    `gorm:"column:username" json:"username"`
	Email    string    `gorm:"column:email" json:"email"`
	Pass     string    `gorm:"column:password" json:"password"`
	Active   time.Time `gorm:"column:activity" json:"activity"`
	SignUp   time.Time `gorm:"column:signup" json:"signup"`
	LoggedIn bool      `gorm:"column:logged" json:"logged"`
}

type Group struct {
	ID      int       `gorm:"primaryKey"`
	Name    string    `gorm:"column:name" json:"name"`
	Desc    string    `gorm:"column:desc" json:"desc"`
	Created time.Time `gorm:"column:created" json:"created"`
}

type Message struct {
	ID     int       `gorm:"primaryKey"`
	Posted time.Time `gorm:"column:posted" json:"posted"`
}

type Priv struct {
	Add_user    bool
	Delete_user bool
	Settings    bool // title, desc, nicks
}

type Member struct {
	id_group int    `gorm:"primaryKey;column:id_group"`
	id_user  int    `gorm:"primaryKey;column:id_user"`
	nick     string `gorm:"column:nick"`
	priv     Priv   `gorm:"embedded"`
}
