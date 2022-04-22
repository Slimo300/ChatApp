package orm

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"gorm.io/gorm"
)

func (db *Database) SendGroupInvite(issID, group uint, targetName string) (models.Invite, error) {

	// Finding new user by username
	var target models.User
	if err := db.Where(&models.User{UserName: targetName}).First(&target).Error; err != nil {
		return models.Invite{}, err
	}
	// Creating invite
	invite := models.Invite{IssId: issID, TargetID: target.ID, GroupID: group, Created: time.Now(), Modified: time.Now(), Status: database.INVITE_AWAITING}
	// Saving invite to database
	if err := db.Create(&invite).Error; err != nil {
		return models.Invite{}, err
	}

	return invite, nil
}

func (db *Database) RespondGroupInvite(userID, inviteID uint, response bool) (models.Group, error) {

	var group models.Group
	var invite models.Invite

	if err := db.First(&invite, inviteID).Error; err != nil {
		return models.Group{}, nil
	}
	// Checking if user is the receiver of invite
	if invite.TargetID != userID {
		return models.Group{}, database.ErrNoPrivilages
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

func (db *Database) SendFriendsInvite(issuerID uint, targetName string) (models.FriendsInvite, error) {

	var target models.User

	if err := db.Where(models.User{UserName: targetName}).First(&target).Error; err != nil {
		return models.FriendsInvite{}, err
	}

	invite := models.FriendsInvite{IssId: issuerID, TargetID: target.ID, Status: database.INVITE_AWAITING, Created: time.Now(), Modified: time.Now()}

	if err := db.Create(&invite).Error; err != nil {
		return models.FriendsInvite{}, err
	}

	return invite, nil
}

func (db *Database) RespondFriendsInvite(userID, inviteID uint, response bool) (models.Group, error) {

	var invite models.FriendsInvite
	if err := db.First(&invite, inviteID).Error; err != nil {
		return models.Group{}, nil
	}
	// checking whether responding user is in fact the target of invitation
	if invite.TargetID != userID {
		return models.Group{}, database.ErrNoPrivilages
	}

	var group models.Group
	if response { // user accepts the invite
		if err := db.Transaction(func(tx *gorm.DB) error {
			group = models.Group{Name: "", Desc: "", Created: time.Now()}
			if err := db.Create(&group).Error; err != nil {
				return err
			}
			// Creating memberships for both users both get only creator priv to delete group, no adding, deleting or setting is allowed
			if err := db.createMemberFromID(invite.IssId, group.ID, false, false, false, true); err != nil {
				return err
			}
			if err := db.createMemberFromID(invite.TargetID, group.ID, false, false, false, true); err != nil {
				return err
			}
			if err := db.Model(&invite).Updates(models.FriendsInvite{Status: database.INVITE_ACCEPT, Modified: time.Now()}).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return models.Group{}, err
		}
	} else { // user declines the invite
		if err := db.Model(&invite).Updates(models.FriendsInvite{Status: database.INVITE_DECLINE, Modified: time.Now()}).Error; err != nil {
			return models.Group{}, err
		}
	}

	return group, nil
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
