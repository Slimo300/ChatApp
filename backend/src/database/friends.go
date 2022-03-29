package database

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

func (db *Database) AddFriend(id int, username string) (models.Group, error) {

	var issuer models.User
	if err := db.Where(models.User{ID: uint(id)}).First(&issuer).Error; err != nil {
		return models.Group{}, err
	}

	var user models.User
	if err := db.Where(models.User{UserName: username}).First(&user).Error; err != nil {
		return models.Group{}, err
	}

	group := models.Group{Name: "", Desc: "", Created: time.Now()}
	if err := db.Create(&group).Error; err != nil {
		return models.Group{}, nil
	}
	member1 := models.Member{GroupID: group.ID, UserID: issuer.ID, Nick: issuer.UserName, Adding: false, Deleting: false, Setting: false, Creator: true}
	member2 := models.Member{GroupID: group.ID, UserID: user.ID, Nick: user.UserName, Adding: false, Deleting: false, Setting: false, Creator: true}
	if err := db.Create(member1).Error; err != nil {
		return models.Group{}, nil
	}
	if err := db.Create(member2).Error; err != nil {
		return models.Group{}, nil
	}

	return models.Group{}, nil
}
