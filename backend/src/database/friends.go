package database

import "github.com/Slimo300/ChatApp/backend/src/models"

func (db *Database) AddFriend(id int, username string) (models.Invite, error) {

	var issuer models.User
	if err := db.Where(models.User{ID: uint(id)}).First(&issuer).Error; err != nil {
		return models.Invite{}, err
	}

	var user models.User
	if err := db.Where(models.User{UserName: username}).First(&user).Error; err != nil {
		return models.Invite{}, err
	}

	invite := models.Invite{IssId: issuer.ID, TargetID: user.ID, Status: INVITE_AWAITING}
	if err := db.Create(&invite).Error; err != nil {
		return models.Invite{}, err
	}

	return invite, nil
}
