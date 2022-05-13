package orm

import (
	"errors"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"gorm.io/gorm"
)

// GetUserGroups returns a slice of Groups of which user is a member
func (db *Database) GetUserGroups(id uint) (groups []models.Group, err error) {

	// Getting ids of groups in which user has membership
	var usergroups []uint
	if err := db.Table("`groups`").Select("`groups`.id").
		Joins("inner join `members` on `members`.group_id = `groups`.id").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Where("`members`.deleted = false").
		Where("`users`.id = ?", id).Scan(&usergroups).Error; err != nil {
		return groups, err
	}

	// Getting full groups data with its members
	if err := db.Where("id in (?)", usergroups).Preload("Members", "deleted is false").Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
}

func (db *Database) CreateGroup(id uint, name, desc string) (models.Group, error) {
	group := models.Group{Name: name, Desc: desc, Created: time.Now()}

	var issuer models.User
	db.Where(models.User{ID: id}).First(&issuer)

	if err := db.Transaction(func(tx *gorm.DB) error {
		creation := tx.Create(&group)
		if creation.Error != nil {
			return creation.Error
		}
		member := models.Member{UserID: id, GroupID: group.ID, Adding: true, Deleting: true, Setting: true, Creator: true, Nick: issuer.UserName}
		m_create := tx.Create(&member)
		if m_create.Error != nil {
			return m_create.Error
		}

		return nil
	}); err != nil {
		return models.Group{}, err
	}

	if err := db.Where(models.Group{ID: group.ID}).Preload("Members", "deleted is false").First(&group).Error; err != nil {
		return models.Group{}, err
	}
	return group, nil
}

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

// Deletes a specified group if user is authorized to do so
func (db *Database) DeleteGroup(id_group, id_user uint) (models.Group, error) {
	// getting user membership to check his privilages
	var membership models.Member
	if err := db.Where(&models.Member{UserID: id_user, GroupID: id_group}).First(&membership).Error; err != nil {
		return models.Group{}, err
	}

	// checking whether user have privilages to delete a group
	if !membership.Creator {
		return models.Group{}, database.ErrNoPrivilages
	}

	// deleting memberships
	if err := db.Where(models.Member{GroupID: id_group}).Delete(&models.Member{}).Error; err != nil {
		return models.Group{}, err
	}
	// deleting specified group
	group := models.Group{ID: id_group}
	if err := db.Delete(&group).Error; err != nil {
		return models.Group{}, err
	}

	return group, nil
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
