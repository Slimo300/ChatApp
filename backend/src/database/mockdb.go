package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

type MockDB struct {
	Users []models.User
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
		}
	]
	`

	var users []models.User
	json.Unmarshal([]byte(USERS), &users)

	// add data
	return &MockDB{Users: users}
}

// GetUserById(id int) (models.User, error)
// RegisterUser(models.User) (models.User, error)

// SignInUser(name string, pass string) (models.User, error)
// SignOutUser(email string) error

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
