package mock

import (
	"errors"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"gorm.io/gorm"
)

func (m *MockDB) DeleteUserFromGroup(memberID uint) (models.Member, error) {

	for _, member := range m.Members {
		if member.ID == memberID {
			member.Deleted = false
			return member, nil
		}
	}
	return models.Member{}, errors.New("no member")

}

func (m *MockDB) GetMemberByID(memberID uint) (models.Member, error) {
	for _, member := range m.Members {
		if member.ID == memberID {
			return member, nil
		}
	}
	return models.Member{}, gorm.ErrRecordNotFound
}

func (m *MockDB) GrantPriv(memberID uint, adding, deleting, setting bool) error {
	for _, member := range m.Members {
		if member.ID == memberID {
			member.Adding = adding
			member.Deleting = deleting
			member.Setting = setting
			return nil
		}
	}

	return errors.New("internal error")
}

func (m *MockDB) IsUserInGroup(userID, groupID uint) bool {
	for _, member := range m.Members {
		if member.GroupID == groupID && member.UserID == userID && !member.Deleted {
			return true
		}
	}
	return false
}
