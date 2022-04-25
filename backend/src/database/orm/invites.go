package orm

import (
	"errors"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"gorm.io/gorm"
)

func (db *Database) SendGroupInvite(issID, group uint, targetName string) (invite models.Invite, err error) {

	// checking whether issuer has rights to add to group
	if err = db.Where(models.Member{GroupID: group, UserID: issID, Adding: true}).First(&models.Member{}).Error; err != nil {
		err = database.ErrNoPrivilages
		return
	}
	// Finding new user by username
	var target models.User
	if err = db.Where(&models.User{UserName: targetName}).First(&target).Error; err != nil {
		return
	}
	// checking if target isn't already in a group
	if err = db.Where(&models.Member{GroupID: group, UserID: target.ID}).First(models.Member{}).Error; err != nil {
		err = errors.New("user already in a group")
		return
	}
	// checking if target isn't already invited
	if err = db.Where(&models.Invite{GroupID: group, TargetID: target.ID}).Error; err != nil {
		err = errors.New("invite already sent")
		return
	}
	// Creating invite
	invite = models.Invite{IssId: issID, TargetID: target.ID, GroupID: group, Created: time.Now(), Modified: time.Now(), Status: database.INVITE_AWAITING}
	// Saving invite to database
	if err = db.Create(&invite).Error; err != nil {
		return
	}

	return invite, nil
}

func (db *Database) RespondGroupInvite(userID, inviteID uint, response bool) (models.Group, error) {

	var group models.Group
	var invite models.Invite

	if err := db.Where(models.Invite{ID: inviteID, TargetID: userID}).First(&invite).Error; err != nil {
		return models.Group{}, nil
	}

	if response { // user joins the group
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := db.createMemberFromID(userID, invite.GroupID, false, false, false, false); err != nil {
				return err
			}
			// Finding group to return
			if err := db.First(&group, invite.GroupID).Error; err != nil {
				return err
			}
			// Updating invite
			if err := db.Model(&invite).Updates(models.Invite{Status: database.INVITE_ACCEPT, Modified: time.Now()}).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return models.Group{}, err
		}
	} else { // user declines invite
		if err := db.Model(&invite).Updates(models.Invite{Status: database.INVITE_DECLINE, Modified: time.Now()}).Error; err != nil {
			return models.Group{}, err
		}
	}

	return group, nil
}

func (db *Database) GetUserInvites(userID uint) (invites []models.Invite, err error) {

	err = db.Where(models.Invite{TargetID: userID, Status: database.INVITE_AWAITING}).Find(&invites).Error
	return
}

// helper for creating membership with id, it first find user to get his
// username and use it as member's nick
func (db *Database) createMemberFromID(userID, groupID uint, adding, deleting, setting, creator bool) error {

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}

	member := models.Member{GroupID: groupID, UserID: userID, Nick: user.UserName, Adding: adding, Deleting: deleting, Setting: setting, Creator: creator}
	if err := db.Create(&member).Error; err != nil {
		return err
	}

	return nil
}
