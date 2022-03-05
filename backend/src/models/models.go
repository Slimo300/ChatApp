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
	Members  []Member  `gorm:"foreignKey:UserID"`
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
	MemberID uint      `gorm:"column:id_member"`
	Member   Member
}

func (Message) TableName() string {
	return "messages"
}

type Priv struct {
	ID       uint `gorm:"primaryKey"` // title, desc, nicks
	Adding   bool `gorm:"column:adding"`
	Deleting bool `gorm:"column:deleting"`
	Setting  bool `gorm:"column:setting"`
	Creator  bool `gorm:"column:creator"`
}

func (Priv) TableName() string {
	return "priv"
}

type Member struct {
	ID      uint   `gorm:"primaryKey"`
	GroupID uint   `gorm:"column:group_id"`
	UserID  uint   `gorm:"column:user_id"`
	Nick    string `gorm:"column:nick"`
	PrivID  uint   `gorm:"column:id_priv"`
	Priv    Priv   `gorm:"foreignKey:ID"`
}

func (Member) TableName() string {
	return "members"
}
