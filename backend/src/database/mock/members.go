package mock

import (
	"errors"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

// Adding user to a group
func (m *MockDB) AddUserToGroup(username string, id_group uint, id_user uint) (models.Member, error) {

	var added models.User // user who is added by his username

	// finding issuer and added
	for _, user := range m.Users {
		if user.UserName == username {
			added = user
		}
	}
	if added.ID == 0 {
		return models.Member{}, errors.New("row not found")
	}

	var membership models.Member
	// getting issuer membership
	for _, mem := range m.Members {
		if mem.UserID == id_user && mem.GroupID == id_group {
			membership = mem
		}
	}
	if (!membership.Creator && !membership.Adding) || membership.ID == 0 {
		return models.Member{}, database.ErrNoPrivilages
	}

	member := models.Member{ID: uint(len(m.Members) + 1), GroupID: id_group, UserID: added.ID, Nick: username, Adding: false,
		Deleting: false, Setting: false, Creator: false, Deleted: false}

	m.Members = append(m.Members, member)

	return member, nil
}

func (m *MockDB) DeleteUserFromGroup(id_member, id_user uint) (models.Member, error) {

	// Getting member to be deleted
	var member *models.Member
	for i, mem := range m.Members {
		if mem.ID == id_member {
			member = &m.Members[i]
		}
	}
	if member == nil {
		return models.Member{}, errors.New("row not found")
	}
	// Checking issuer privilages
	var issuer models.Member
	for _, mem := range m.Members {
		if mem.UserID == id_user && mem.GroupID == member.GroupID {
			issuer = mem
		}
	}
	if issuer.ID == 0 {
		return models.Member{}, errors.New("row not found")
	}
	if !issuer.Deleting {
		return models.Member{}, database.ErrNoPrivilages
	}

	member.Deleted = true

	return *member, nil
}

func (m *MockDB) GrantPriv(id_mem, id uint, adding, deleting, setting bool) error {
	// getting member to be changed
	var member *models.Member
	for i, mem := range m.Members {
		if mem.ID == id_mem {
			member = &m.Members[i]
		}
	}
	if member == nil {
		return errors.New("row not found")
	}
	if member.Deleted {
		return errors.New("member deleted")
	}
	if member.Creator {
		return errors.New("creator can't be modified")
	}
	// getting issuer
	var issuer models.Member
	for _, mem := range m.Members {
		if mem.GroupID == member.GroupID && mem.UserID == id {
			issuer = mem
		}
	}
	// returning errors
	if issuer.ID == 0 {
		return errors.New("row not found")
	}
	if !issuer.Setting {
		return database.ErrNoPrivilages
	}
	// making changes via pointer
	member.Adding = adding
	member.Deleting = deleting
	member.Setting = setting

	return nil
}
