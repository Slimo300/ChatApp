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

func (User) TableName() string {
	return "users"
}

type Group struct {
	ID      int       `gorm:"primaryKey"`
	Name    string    `gorm:"column:name" json:"name"`
	Desc    string    `gorm:"column:desc" json:"desc"`
	Created time.Time `gorm:"column:created" json:"created"`
}

func (Group) TableName() string {
	return "user_groups"
}

type Message struct {
	ID     int       `gorm:"primaryKey"`
	Posted time.Time `gorm:"column:posted" json:"posted"`
}

func (Message) TableName() string {
	return "messages"
}

type Priv struct {
	ID int `gorm:"primaryKey;column:id"` // title, desc, nicks
}

func (Priv) TableName() string {
	return "rights"
}

type Member struct {
	id_group int    `gorm:"primaryKey;column:id_group"`
	id_user  int    `gorm:"primaryKey;column:id_user"`
	nick     string `gorm:"column:nick"`
	priv     Priv   `gorm:"embedded"`
}

func (Member) TableName() string {
	return "members"
}
