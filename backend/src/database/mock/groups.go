package mock

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
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

func (m *MockDB) DeleteGroup(id_group, id_user uint) (models.Group, error) {

	var issuer models.Member
	for _, member := range m.Members {
		if member.GroupID == id_group && member.UserID == id_user {
			issuer = member
			break
		}
	}
	if !issuer.Creator {
		return models.Group{}, database.ErrNoPrivilages
	}
	var deleted_group models.Group
	for i, group := range m.Groups {
		if group.ID == id_group {
			deleted_group = group
			m.Groups = append(m.Groups[:i], m.Groups[i+1:]...)
			break
		}
	}
	return deleted_group, nil
}
