package models

import (
	"time"
)

type User struct {
	ID       uint      `gorm:"primaryKey"`
	UserName string    `gorm:"column:username;unique" json:"username"`
	Email    string    `gorm:"column:email;unique" json:"email"`
	Pass     string    `gorm:"column:password" json:"password"`
	Active   time.Time `gorm:"column:activity" json:"activity"`
	SignUp   time.Time `gorm:"column:signup" json:"signup"`
	LoggedIn bool      `gorm:"column:logged" json:"logged"`
}

func (User) TableName() string {
	return "users"
}

type Group struct {
	ID      uint      `gorm:"primaryKey"`
	Name    string    `gorm:"column:name" json:"name"`
	Desc    string    `gorm:"column:desc" json:"desc"`
	Created time.Time `gorm:"column:created" json:"created"`
}

func (Group) TableName() string {
	return "groups"
}

type Message struct {
	ID       uint      `gorm:"primaryKey"`
	Posted   time.Time `gorm:"column:posted" json:"posted"`
	Text     string    `gorm:"column:text" json:"text"`
	MemberID uint      `gorm:"column:id_member" json:"member_id"`
	Member   Member
}

func (Message) TableName() string {
	return "messages"
}

type Member struct {
	ID       uint `gorm:"primaryKey"`
	GroupID  uint `gorm:"column:group_id"`
	UserID   uint `gorm:"column:user_id"`
	User     User
	Group    Group
	Nick     string `gorm:"column:nick"`
	Adding   bool   `gorm:"column:adding"`
	Deleting bool   `gorm:"column:deleting"`
	Setting  bool   `gorm:"column:setting"`
	Creator  bool   `gorm:"column:creator"`
}

func (Member) TableName() string {
	return "members"
}
