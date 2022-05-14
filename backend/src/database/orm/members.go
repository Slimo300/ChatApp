package orm

import (
	"github.com/Slimo300/ChatApp/backend/src/models"
)

func (db *Database) GetMemberByID(memberID uint) (member models.Member, err error) {
	return member, db.First(&member, memberID).Error
}

func (db *Database) DeleteUserFromGroup(memberID uint) (member models.Member, err error) {
	return member, db.First(&member, memberID).Update("deleted", true).Error
}

func (db *Database) GrantPriv(memberID uint, adding, deleting, setting bool) error {
	return db.First(&models.Member{}, memberID).Updates(models.Member{Adding: adding, Deleting: deleting, Setting: setting}).Error
}
