package orm

import (
	"errors"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"gorm.io/gorm"
)

func (db *Database) AddUserToGroup(username string, id_group uint, id_user uint) (models.Member, error) {

	var issuer models.Member
	if err := db.Table("`members`").Select("members.*").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Joins("inner join `groups` on `groups`.id = `members`.group_id").
		Where("`users`.id = ?", id_user).
		Where("`groups`.id = ?", id_group).Scan(&issuer).Error; err != nil {
		return models.Member{}, err
	}

	if !issuer.Adding {
		return models.Member{}, database.ErrNoPrivilages
	}

	var user models.User
	if err := db.Where(&models.User{UserName: username}).First(&user).Error; err != nil {
		return models.Member{}, err
	}

	var member models.Member
	if err := db.Where(&models.Member{UserID: user.ID, GroupID: id_group}).First(&member).Error; err != nil && err != gorm.ErrRecordNotFound {
		return models.Member{}, err
	}
	if member.Deleted == true {
		if err := db.Model(&member).Update("deleted", false).Error; err != nil {
			return models.Member{}, err
		}
		return member, nil
	}

	member = models.Member{UserID: user.ID, Nick: user.UserName, GroupID: id_group, Adding: false, Deleting: false, Setting: false, Creator: false}
	if err := db.Create(&member).Error; err != nil {
		return models.Member{}, err
	}

	return member, nil
}

func (db *Database) DeleteUserFromGroup(id_member, id_user uint) (models.Member, error) {

	var deleted_member models.Member
	if err := db.First(&deleted_member, id_member).Error; err != nil {
		return models.Member{}, err
	}

	var issuer_member models.Member
	if err := db.Table("`members`").Select("`members`.*").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Joins("inner join `groups` on `groups`.id = `members`.group_id").
		Where("`users`.id = ?", id_user).
		Where("`groups`.id = ?", deleted_member.GroupID).Scan(&issuer_member).Error; err != nil {
		return models.Member{}, err
	}

	if !issuer_member.Deleting {
		return models.Member{}, database.ErrNoPrivilages
	}

	if err := db.Model(&deleted_member).Update("deleted", true).Error; err != nil {
		return models.Member{}, err
	}

	return deleted_member, nil
}

// Updates user rights to a group
func (db *Database) GrantPriv(id_mem, id uint, adding, deleting, setting bool) error {

	var member models.Member
	if err := db.First(&member, id_mem).Error; err != nil {
		return err
	}
	if member.Deleted {
		return errors.New("member deleted")
	}
	if member.Creator {
		return errors.New("creator can't be modified")
	}

	// checking membership of an issuer
	var issuer models.Member
	if err := db.Where(models.Member{UserID: id, GroupID: member.GroupID}).First(&issuer).Error; err != nil {
		return err
	}
	if !issuer.Setting {
		return database.ErrNoPrivilages
	}

	if err := db.Model(member).Updates(models.Member{Adding: adding, Deleting: deleting, Setting: setting}).Error; err != nil {
		return err
	}

	return nil
}
