package mock

import (
	"errors"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

func (m *MockDB) SendGroupInvite(issID, group uint, targetName string) (models.Invite, error) {
	ErrFlag := false
	for _, mem := range m.Members {
		if mem.GroupID == group && mem.UserID == issID && mem.Adding {
			ErrFlag = true
			break
		}
	}
	if !ErrFlag {
		return models.Invite{}, database.ErrNoPrivilages
	}

	var user models.User
	ErrFlag = false
	for i, u := range m.Users {
		if u.UserName == targetName {
			user = m.Users[i]
			ErrFlag = true
		}
	}
	if !ErrFlag {
		return models.Invite{}, errors.New("user not found")
	}

	// checking if user is not already in a group
	for _, mem := range m.Members {
		if mem.UserID == user.ID && mem.GroupID == group {
			return models.Invite{}, errors.New("user already in a group")
		}
	}

	// checking if invite is not a duplicate
	for _, inv := range m.Invites {
		if inv.GroupID == group && inv.TargetID == user.ID && inv.Status == database.INVITE_AWAITING {
			return models.Invite{}, errors.New("invite already sent")
		}
	}

	invite := models.Invite{ID: uint(len(m.Invites) + 1), IssId: issID, TargetID: user.ID, GroupID: group, Status: 0, Created: time.Now(), Modified: time.Now()}
	m.Invites = append(m.Invites, invite)
	return invite, nil

}

func (mock *MockDB) GetUserInvites(userID uint) ([]models.Invite, error) {

	userInvites := []models.Invite{}

	for _, invite := range mock.Invites {
		if invite.TargetID == userID && invite.Status == database.INVITE_AWAITING {
			userInvites = append(userInvites, invite)
		}
	}

	return userInvites, nil
}

func (mock *MockDB) RespondGroupInvite(userID, inviteID uint, response bool) (models.Group, error) {

	var respondedInvite *models.Invite
	for i, invite := range mock.Invites {
		if invite.ID == inviteID && invite.TargetID == userID {
			respondedInvite = &mock.Invites[i]
		}
	}
	if respondedInvite == nil {
		return models.Group{}, errors.New("no such invite")
	}

	if response {
		var respondingUser *models.User
		for i, user := range mock.Users {
			if user.ID == userID {
				respondingUser = &mock.Users[i] // no need to check if assigned in mock database
			}
		}

		mock.Members = append(mock.Members, models.Member{
			ID:       uint(len(mock.Members) + 1),
			UserID:   userID,
			GroupID:  respondedInvite.GroupID,
			Nick:     respondingUser.UserName,
			Adding:   false,
			Deleting: false,
			Setting:  false,
			Creator:  false,
			Deleted:  false,
		})

		respondedInvite.Status = database.INVITE_ACCEPT
		respondedInvite.Modified = time.Now()

		for _, group := range mock.Groups {
			if group.ID == respondedInvite.GroupID {
				return group, nil
			}
		}
	} else {
		respondedInvite.Status = database.INVITE_DECLINE
		respondedInvite.Modified = time.Now()

		return models.Group{}, nil
	}

	return models.Group{}, errors.New("Something went wrong")
}
