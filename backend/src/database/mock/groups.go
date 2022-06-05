package mock

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"gorm.io/gorm"
)

func (m *MockDB) CreateGroup(id uint, name, desc string) (models.Group, error) {
	newGroup := models.Group{
		ID:      uint(len(m.Groups) + 1),
		Name:    name,
		Desc:    desc,
		Created: time.Now(),
	}
	m.Groups = append(m.Groups, newGroup)

	return newGroup, nil
}

func (m *MockDB) GetUserGroups(id uint) ([]models.Group, error) {
	var groups []models.Group
	for _, member := range m.Members {
		if member.UserID == id {
			for _, group := range m.Groups {
				if member.GroupID == group.ID {
					groups = append(groups, group)
				}
			}
		}
	}
	return groups, nil
}

func (m *MockDB) DeleteGroup(groupID uint) (models.Group, error) {
	var deletedGroup models.Group
	for i, group := range m.Groups {
		if group.ID == groupID {
			deletedGroup = group
			m.Groups = append(m.Groups[:i], m.Groups[i+1:]...)
			break
		}
	}
	return deletedGroup, nil
}

func (m *MockDB) GetUserGroupMember(userID, groupID uint) (models.Member, error) {
	for _, member := range m.Members {
		if member.GroupID == groupID && member.UserID == userID {
			return member, nil
		}
	}
	return models.Member{}, gorm.ErrRecordNotFound
}

func (m *MockDB) DeleteGroupProfilePicture(groupID uint) error {
	for _, group := range m.Groups {
		if group.ID == groupID {
			group.Picture = ""
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockDB) SetGroupProfilePicture(groupID uint, newURI string) error {
	for _, group := range m.Groups {
		if group.ID == groupID {
			group.Picture = newURI
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockDB) GetGroupProfilePicture(groupID uint) (string, error) {
	for _, group := range m.Groups {
		if group.ID == groupID {
			return group.Picture, nil
		}
	}
	return "", gorm.ErrRecordNotFound
}
