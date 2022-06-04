package models

import (
	"time"
)

type User struct {
	ID       uint      `gorm:"primaryKey"`
	UserName string    `gorm:"column:username;unique" json:"username"`
	Email    string    `gorm:"column:email;unique" json:"email"`
	Pass     string    `gorm:"column:password" json:"password"`
	Picture  string    `gorm:"picture_url" json:"pictureUrl"`
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
	Picture string    `gorm:"column:picture_url" json:"pictureUrl"`
	Created time.Time `gorm:"column:created" json:"created"`
	Members []Member  `gorm:"foreignKey:GroupID"`
}

func (Group) TableName() string {
	return "groups"
}

type Message struct {
	ID       uint      `gorm:"primaryKey"`
	Posted   time.Time `gorm:"column:posted" json:"posted"`
	Text     string    `gorm:"column:text" json:"text"`
	MemberID uint      `gorm:"column:id_member" json:"member_id"`
	Member   Member    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Message) TableName() string {
	return "messages"
}

type Member struct {
	ID       uint   `gorm:"primaryKey"`
	GroupID  uint   `gorm:"column:group_id;uniqueIndex:idx_first" json:"group_id"`
	UserID   uint   `gorm:"column:user_id;uniqueIndex:idx_first" json:"user_id"`
	User     User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group    Group  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Nick     string `gorm:"column:nick" json:"nick"`
	Adding   bool   `gorm:"column:adding" json:"adding"`
	Deleting bool   `gorm:"column:deleting" json:"deleting"`
	Setting  bool   `gorm:"column:setting" json:"setting" `
	Creator  bool   `gorm:"column:creator" json:"creator"`
	Deleted  bool   `gorm:"column:deleted" json:"deleted"`
}

func (Member) TableName() string {
	return "members"
}

type Invite struct {
	ID       uint      `gorm:"primaryKey"`
	IssId    uint      `gorm:"column:iss_id" json:"issID"`
	Iss      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	TargetID uint      `gorm:"column:target_id" json:"targetID"`
	Target   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	GroupID  uint      `gorm:"column:group_id" json:"groupID"`
	Group    Group     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Status   int       `gorm:"column:status" json:"status"`
	Created  time.Time `gorm:"column:created" json:"created"`
	Modified time.Time `gorm:"column:modified" json:"modified"`
}

func (Invite) TableName() string {
	return "invites"
}
