package database

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

type MockDB struct {
	users    []models.User
	groups   []models.Group
	messages []models.Message
	members  []models.Member
}

func NewMockDB() *MockDB {

	USERS := `[
		{
				"ID": 1,
				"signup": "2018-08-14T07:52:54Z",
				"active": "2019-01-13T22:00:45Z",
				"username": "Mal",
				"email": "mal.zein@email.com",
				"password": "$2a$10$BUOsZ9O5Lt/YJv31gMio/.HvuOTiu7waiE936g7dnKQ37VY8he2GW",
				"logged": true,
		},
		{
				"ID": 2,
				"signup": "2018-08-14T07:52:55Z",
				"active": "2019-01-12T22:39:01Z",
				"username": "River",
				"email": "river.sam@email.com",
				"password": "$2a$10$mNbCLmfCAc0.4crDg3V3fe0iO1yr03aRfE7Rr3vdfKMGVnnzovCZq",
				"logged": false,
		},
		{
				"ID": 3,
				"signup": "2019-01-13T08:43:44Z",
				"active": "2019-01-13T15:12:25Z",
				"username": "John",
				"email": "john.doe@bla.com",
				"password": "$2a$10$T4c8rmpbgKrUA0sIqtHCaO0g2XGWWxFY4IGWkkpVQOD/iuBrwKrZu",
				"logged": false,
		}
	]
	`

	var users []models.User
	json.Unmarshal([]byte(USERS), &users)

	// add data
	return &MockDB{users: users}
}

// GetUserById(id int) (models.User, error)
// RegisterUser(models.User) (models.User, error)

// SignInUser(name string, pass string) (models.User, error)
// SignOutUser(email string) error

func (m *MockDB) GetUserById(id int) (models.User, error) {
	for _, user := range m.users {
		if user.ID == uint(id) {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("No user with id: %d", id)
}

func (m *MockDB) RegisterUser(user models.User) (models.User, error) {
	user.ID = uint(len(m.users) + 1)
	user.Active = time.Now()
	user.SignUp = time.Now()
	user.LoggedIn = false
	m.users = append(m.users, user)
	return user, nil
}

func (m *MockDB) SignInUser(name, pass string) (models.User, error) {
	for _, user := range m.users {
		if !strings.EqualFold(user.Email, name) {
			continue
		}
		if checkPassword(user.Pass, pass) {
			user.LoggedIn = true
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("No email %s in database", name)
}

func (m *MockDB) SignOutUser(email string) error {
	for _, user := range m.users {
		if !strings.EqualFold(user.Email, email) {
			continue
		}
		user.LoggedIn = false
		return nil
	}
	return fmt.Errorf("No user with email: %s", email)
}
