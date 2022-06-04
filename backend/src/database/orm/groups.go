package orm

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"gorm.io/gorm"
)

func (db *Database) GetUserGroups(id uint) (groups []models.Group, err error) {

	var userGroupsIDs []uint
	if err := db.Table("`groups`").Select("`groups`.id").
		Joins("inner join `members` on `members`.group_id = `groups`.id").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Where("`members`.deleted = false").
		Where("`users`.id = ?", id).Scan(&userGroupsIDs).Error; err != nil {
		return groups, err
	}

	if err := db.Where("id in (?)", userGroupsIDs).Preload("Members", "deleted is false").Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
}

func (db *Database) CreateGroup(id uint, name, desc string) (models.Group, error) {
	group := models.Group{Name: name, Desc: desc, Created: time.Now(), Picture: ""}

	var creator models.User
	db.Where(models.User{ID: id}).First(&creator)

	if err := db.Transaction(func(tx *gorm.DB) error {
		creation := tx.Create(&group)
		if creation.Error != nil {
			return creation.Error
		}
		member := models.Member{UserID: id, GroupID: group.ID, Adding: true, Deleting: true, Setting: true, Creator: true, Nick: creator.UserName}
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

func (db *Database) DeleteGroup(groupID uint) (group models.Group, err error) {

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Where(models.Member{GroupID: groupID}).Delete(&models.Member{}).Error; err != nil {
			return err
		}
		group = models.Group{ID: groupID}
		if err := db.Delete(&group).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return models.Group{}, err
	}

	return group, nil
}
