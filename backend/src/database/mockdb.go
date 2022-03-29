package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

type MockDB struct {
	Users    []models.User
	Groups   []models.Group
	Members  []models.Member
	Messages []models.Message
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

	var users []models.User
	json.Unmarshal([]byte(USERS), &users)

	var groups []models.Group
	json.Unmarshal([]byte(GROUPS), &groups)

	var members []models.Member
	json.Unmarshal([]byte(MEMBERS), &members)

	var messages []models.Message
	json.Unmarshal([]byte(MESSAGES), &messages)

	// add data
	return &MockDB{Users: users, Groups: groups, Members: members, Messages: messages}
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
	pass, err := hashPassword(user.Pass)
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
		if checkPassword(user.Pass, pass) {
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

func (m *MockDB) GetGroupMessages(id_user, id_group, offset uint) ([]Message, error) {
	var messages []Message

	for _, member := range m.Members {
		if member.GroupID == id_group && member.UserID == id_user {
			break
		}
		return nil, errors.New("User cannot request from this group")
	}

	for _, message := range m.Messages {
		for _, member := range m.Members {
			if message.MemberID == member.ID && member.GroupID == id_group {
				messages = append(messages, Message{
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

func (m *MockDB) GetGroupMembership(group, user uint) (models.Member, error) {
	for _, member := range m.Members {
		if member.GroupID == group && member.UserID == user {
			return member, nil
		}
	}
	return models.Member{}, errors.New("Err no record")
}

func (m *MockDB) AddMessage(msg Message) (Message, error) {
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

func (m *MockDB) DeleteGroup(id_group, id_user uint) error {
	groupDel := false
	for i, group := range m.Groups {
		if group.ID == id_group {
			for _, member := range m.Members {
				if member.GroupID == id_group && member.UserID == id_user && member.Creator {
					m.Groups = append(m.Groups[:i], m.Groups[i+1:]...)
					groupDel = true
					break
				}
				return errors.New("Couldn't delete group")
			}
			break
		}
	}
	if groupDel {
		var newMembers []models.Member
		for _, member := range m.Members {
			if member.GroupID != id_group {
				newMembers = append(newMembers, member)
			}
		}
		m.Members = newMembers
	}
	return nil
}

// func AddFriend takes "id" which is id of issuer and username of invited user
func (m *MockDB) AddFriend(id int, username string) (models.Invite, error) {
	return models.Invite{}, nil
}

// RespondInvite takes invite id and response which is 1 (agree) or 2 (decline)
func (m *MockDB) RespondInvite(id_inv, response int) (models.Group, error) {
	return models.Group{}, nil
}

// Adding user to a group
func (m *MockDB) AddUserToGroup(username string, id_group uint, id_user uint) error {

	var added models.User // user who is added by his username

	// finding issuer and added
	for _, user := range m.Users {
		if user.UserName == username {
			added = user
		}
	}
	if added.ID == 0 {
		return errors.New("row not found")
	}

	var membership models.Member
	// getting issuer membership
	for _, mem := range m.Members {
		if mem.UserID == id_user && mem.GroupID == id_group {
			membership = mem
		}
	}
	if (!membership.Creator && !membership.Adding) || membership.ID == 0 {
		return ErrNoPrivilages
	}

	m.Members = append(m.Members, models.Member{ID: uint(len(m.Members) + 1), GroupID: id_group, UserID: added.ID, Nick: username, Adding: false,
		Deleting: false, Setting: false, Creator: false, Deleted: false})

	return nil
}

func (m *MockDB) DeleteUserFromGroup(id_member, id_user uint) error {

	// Getting member to be deleted
	var member *models.Member
	for i, mem := range m.Members {
		if mem.ID == id_member {
			member = &m.Members[i]
		}
	}
	if member == nil {
		return errors.New("row not found")
	}
	// Checking issuer privilages
	var issuer models.Member
	for _, mem := range m.Members {
		if mem.UserID == id_user && mem.GroupID == member.GroupID {
			issuer = mem
		}
	}
	if issuer.ID == 0 {
		return errors.New("row not found")
	}
	if !issuer.Deleting && !issuer.Creator {
		return ErrNoPrivilages
	}

	member.Deleted = false

	return nil
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
