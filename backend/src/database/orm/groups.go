package orm

import (
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
