package mock

import (
	"gorm.io/gorm"
)

func (m *MockDB) SetPassword(userID uint, password string) error {
	for _, user := range m.Users {
		if user.ID == userID {
			user.Pass = password
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockDB) DeleteProfilePicture(userID uint) error {
	for _, user := range m.Users {
		if user.ID == userID {
			user.Picture = ""
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockDB) SetProfilePicture(userID uint, newURI string) error {
	for _, user := range m.Users {
		if user.ID == userID {
			user.Picture = newURI
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockDB) GetProfilePictureURL(userID uint) (string, error) {
	for _, user := range m.Users {
		if user.ID == userID {
			return user.Picture, nil
		}
	}
	return "", gorm.ErrRecordNotFound
}
