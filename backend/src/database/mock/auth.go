package mock

import (
	"fmt"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

func (m *MockDB) GetUserById(id int) (models.User, error) {
	for _, user := range m.Users {
		if user.ID == uint(id) {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("No user with id: %d", id)
}

func (m *MockDB) GetUserByEmail(email string) (models.User, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("No user with email: %s", email)
}

func (m *MockDB) RegisterUser(user models.User) (models.User, error) {
	user.ID = uint(len(m.Users) + 1)
	user.Active = time.Now()
	user.SignUp = time.Now()
	user.LoggedIn = false
	m.Users = append(m.Users, user)
	return user, nil
}

func (m *MockDB) SignInUser(id uint) error {
	for _, user := range m.Users {
		if user.ID == id {
			user.LoggedIn = true
			return nil
		}
	}
	return fmt.Errorf("No user with id: %d", id)
}

func (m *MockDB) SignOutUser(id uint) error {
	for _, user := range m.Users {
		if user.ID == id {
			user.LoggedIn = false
			return nil
		}
	}
	return fmt.Errorf("No user with id: %d", id)
}

func (m *MockDB) IsEmailInDatabase(email string) bool {
	for _, user := range m.Users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (m *MockDB) IsUsernameInDatabase(username string) bool {
	for _, user := range m.Users {
		if user.UserName == username {
			return true
		}
	}
	return false
}
