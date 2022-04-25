package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

type MockDB struct {
	Users    []models.User
	Groups   []models.Group
	Members  []models.Member
	Messages []models.Message
	Invites  []models.Invite
}

func NewMockDB() *MockDB {

	USERS := `[
		{
			"ID": 1,
			"signup": "2018-08-14T07:52:54Z",
			"active": "2019-01-13T22:00:45Z",
			"username": "Mal",
			"email": "mal.zein@email.com",
			"password": "$2a$10$6BSuuiaPdRJJF2AygYAfnOGkrKLY2o0wDWbEpebn.9Rk0O95D3hW.",
			"logged": true
		},
		{
			"ID": 2,
			"signup": "2018-08-14T07:52:55Z",
			"active": "2019-01-12T22:39:01Z",
			"username": "River",
			"email": "river.sam@email.com",
			"password": "$2a$10$BvQjoN.PH8FkVPV3ZMBK1O.3LGhrF/RhZ2kM/h9M3jPA1d2lZkL.C",
			"logged": false
		},
		{
			"ID": 3,
			"username": "John",
			"signup": "2019-01-13T08:43:44Z",
			"active": "2019-01-13T15:12:25Z",
			"email": "john.doe@bla.com",
			"password": "$2a$10$T4c8rmpbgKrUA0sIqtHCaO0g2XGWWxFY4IGWkkpVQOD/iuBrwKrZu",
			"logged": false
		},
		{
			"ID": 4,
			"username": "Kal",
			"signup": "2019-01-13T08:53:44Z",
			"active": "2019-01-13T15:52:25Z",
			"email": "kal.doe@bla.com",
			"password": "$2a$10$T4c8rmpbgKrUA0sIqtHCaO0g2XGWWxFY4IGWkkpVQOD/iuBrwKrZu",
			"logged": false
		},
		{
			"ID": 5,
			"username": "Kel",
			"signup": "2019-01-12T08:53:44Z",
			"active": "2019-01-12T15:52:25Z",
			"email": "kel.doa@bla.com",
			"password": "$2a$10$T4c8rmpbgKrUA0sIqtHCaO0g2XGWWxFY4IGWkkpVQOD/iuBrwKrZu",
			"logged": false
		}
	]`

	GROUPS := `[
		{
			"ID": 1,
			"name": "New Group",
			"desc": "totally new group",
			"created": "2019-01-13T08:47:44Z"
		}	
	]`

	MEMBERS := `[
		{
			"ID": 1,
			"group_id": 1,
			"user_id": 1,
			"nick": "Mal",
			"adding": true,
			"deleting": true,
			"setting": true,
			"creator": true,
			"deleted": false
		},
		{
			"ID": 2,
			"group_id": 1,
			"user_id": 2,
			"nick": "River",
			"adding": false,
			"deleting": false,
			"setting": false,
			"creator": false,
			"deleted": false
		},
		{
			"ID": 4,
			"group_id": 1,
			"user_id": 4,
			"nick": "Kal",
			"adding": false,
			"deleting": false,
			"setting": false,
			"creator": false,
			"deleted": true
		}
	]`

	MESSAGES := `[
		{
			"ID": 1,
			"posted": "2019-01-13T22:00:45Z",
			"text": "elo",
			"member_id": 1
		},
		{
			"ID": 2,
			"posted": "2019-01-15T22:00:45Z",
			"text": "siema",
			"member_id": 2
		}, 
		{
			"ID": 3,
			"posted": "2019-01-16T22:00:45Z",
			"text": "elo elo",
			"member_id": 1
		},
		{
			"ID": 4,
			"posted": "2019-01-17T22:00:45Z",
			"text": "siema siema",
			"member_id": 2
		}
	]`

	INVITES := `[
		{
			"ID": 1,
			"issID": 1,
			"targetID": 3,
			"groupID": 1,
			"status": 0,
			"created": "2019-03-17T22:04:45Z",
			"modified": "2019-03-17T22:04:45Z"
		}
	]`

	var users []models.User
	json.Unmarshal([]byte(USERS), &users)

	var groups []models.Group
	json.Unmarshal([]byte(GROUPS), &groups)

	var members []models.Member
	json.Unmarshal([]byte(MEMBERS), &members)

	var messages []models.Message
	json.Unmarshal([]byte(MESSAGES), &messages)

	var invites []models.Invite
	json.Unmarshal([]byte(INVITES), &invites)

	// add data
	return &MockDB{Users: users, Groups: groups, Members: members, Messages: messages, Invites: invites}
}

func (m *MockDB) GetUserById(id int) (models.User, error) {
	for _, user := range m.Users {
		if user.ID == uint(id) {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("No user with id: %d", id)
}

func (m *MockDB) RegisterUser(user models.User) (models.User, error) {
	user.ID = uint(len(m.Users) + 1)
	user.Active = time.Now()
	user.SignUp = time.Now()
	user.LoggedIn = false
	pass, err := HashPassword(user.Pass)
	user.Pass = pass
	if err != nil {
		return models.User{}, errors.New("couldn't register user")
	}
	for _, u := range m.Users {
		if user.Email == u.Email {
			return user, errors.New("email taken")
		}
	}
	m.Users = append(m.Users, user)
	return user, nil
}

func (m *MockDB) SignInUser(name, pass string) (models.User, error) {
	for _, user := range m.Users {
		if !(user.Email == name) {
			continue
		}
		if CheckPassword(user.Pass, pass) {
			user.LoggedIn = true
			return user, nil
		} else {
			return models.User{}, ErrINVALIDPASSWORD
		}
	}
	return models.User{}, fmt.Errorf("No email %s in database", name)
}

func (m *MockDB) SignOutUser(id uint) error {
	for _, user := range m.Users {
		if user.ID != id {
			continue
		}
		user.LoggedIn = false
		return nil
	}
	return fmt.Errorf("No user with id: %d", id)
}

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

func (m *MockDB) GetGroupMessages(id_user, id_group uint, offset, num int) ([]communication.Message, error) {
	var messages []communication.Message

	for _, member := range m.Members {
		if member.GroupID == id_group && member.UserID == id_user {
			break
		}
		return nil, errors.New("User cannot request from this group")
	}

	for _, message := range m.Messages {
		for _, member := range m.Members {
			if message.MemberID == member.ID && member.GroupID == id_group {
				messages = append(messages, communication.Message{
					Group:   uint64(id_group),
					Member:  uint64(member.ID),
					Message: message.Text,
					Nick:    member.Nick,
					When:    message.Posted.Format(TIME_FORMAT),
				})
			}
		}
	}

	return messages, nil
}

func (m *MockDB) AddMessage(msg communication.Message) (communication.Message, error) {
	msgTime := time.Now()
	m.Messages = append(m.Messages, models.Message{
		ID:       uint(len(m.Messages) + 1),
		Posted:   msgTime,
		Text:     msg.Message,
		MemberID: uint(msg.Member),
	})

	msg.When = msgTime.Format(TIME_FORMAT)
	return msg, nil
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
		return models.Group{}, ErrNoPrivilages
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
		return models.Member{}, ErrNoPrivilages
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
		return models.Member{}, ErrNoPrivilages
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
		return ErrNoPrivilages
	}
	// making changes via pointer
	member.Adding = adding
	member.Deleting = deleting
	member.Setting = setting

	return nil
}

func (m *MockDB) SendGroupInvite(issID, group uint, targetName string) (models.Invite, error) {
	ErrFlag := false
	for _, mem := range m.Members {
		if mem.GroupID == group && mem.UserID == issID && mem.Adding {
			ErrFlag = true
			break
		}
	}
	if !ErrFlag {
		return models.Invite{}, ErrNoPrivilages
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
		if inv.GroupID == group && inv.TargetID == user.ID && inv.Status == INVITE_AWAITING {
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
		if invite.TargetID == userID && invite.Status == INVITE_AWAITING {
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

		respondedInvite.Status = INVITE_ACCEPT
		respondedInvite.Modified = time.Now()

		for _, group := range mock.Groups {
			if group.ID == respondedInvite.GroupID {
				return group, nil
			}
		}
	} else {
		respondedInvite.Status = INVITE_DECLINE
		respondedInvite.Modified = time.Now()

		return models.Group{}, nil
	}

	return models.Group{}, errors.New("Something went wrong")
}
