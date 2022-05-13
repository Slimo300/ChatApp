package orm

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

//GetUserById returns a user with specified id
func (db *Database) GetUserById(id int) (user models.User, err error) {
	return user, db.First(&user, id).Error
}

//RegisterUser adds a new user to database
func (db *Database) RegisterUser(user models.User) (models.User, error) {
	pass, err := database.HashPassword(user.Pass)
	if err != nil {
		return models.User{}, err
	}
	user.Pass = pass
	user.Active = time.Now()
	user.SignUp = time.Now()
	user.LoggedIn = false
	return user, db.Create(&user).Error
}

// SignInUser validates user credentials and return his data if they are propper
func (db *Database) SignInUser(email, pass string) (user models.User, err error) {
	result := db.Table("users").Where(&models.User{Email: email})
	err = result.First(&user).Error
	if err != nil {
		return user, err
	}
	if !database.CheckPassword(user.Pass, pass) {
		return user, database.ErrINVALIDPASSWORD
	}
	user.Pass = ""
	err = result.Update("logged", 1).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (db *Database) SignOutUser(id uint) error {
	user := db.Table("users").First(&models.User{ID: id})
	if user.Error != nil {
		return user.Error
	}
	if err := user.Updates(models.User{LoggedIn: false, Active: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}
