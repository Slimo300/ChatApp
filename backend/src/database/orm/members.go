package orm

import (
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/google/uuid"
)

func (db *Database) GetMemberByID(memberID uuid.UUID) (member models.Member, err error) {
	return member, db.First(&member, memberID).Error
}

func (db *Database) DeleteUserFromGroup(memberID uuid.UUID) (member models.Member, err error) {
	return member, db.First(&member, memberID).Update("deleted", true).Error
}

func (db *Database) GrantPriv(memberID uuid.UUID, adding, deleting, setting bool) error {
	return db.First(&models.Member{}, memberID).Updates(models.Member{Adding: adding, Deleting: deleting, Setting: setting}).Error
}

func (db *Database) GetUserGroupMember(userID, groupID uuid.UUID) (member models.Member, err error) {
	return member, db.Where(models.Member{UserID: userID, GroupID: groupID}).First(&member).Error
}

func (db *Database) IsUserInGroup(userID, groupID uuid.UUID) bool {
	var member models.Member
	err := db.Where(models.Member{UserID: userID, GroupID: groupID}).First(&member).Error
	if err != nil || member.Deleted == true {
		return false
	}
	return true
}
