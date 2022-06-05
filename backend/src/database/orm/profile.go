package orm

import "github.com/Slimo300/ChatApp/backend/src/models"

func (db *Database) SetPassword(userID uint, password string) error {
	return db.First(&models.User{}, userID).Update("password", password).Error
}

func (db *Database) DeleteProfilePicture(userID uint) error {
	return db.First(&models.User{}, userID).Update("picture", "").Error
}

func (db *Database) SetProfilePicture(userID uint, newURI string) error {
	return db.First(&models.User{}, userID).Update("picture", newURI).Error
}

func (db *Database) GetProfilePictureURL(userID uint) (string, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return "", err
	}
	return user.Picture, nil
}
