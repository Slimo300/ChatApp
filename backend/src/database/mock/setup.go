package mock

import (
	"encoding/json"

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
