package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/Slimo300/ChatApp/backend/src/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	*gorm.DB
}

// Setup creates Database object and initializes connection between MySQL database
func Setup() (*Database, error) {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@/%s?parseTime=true", os.Getenv("MYSQLUSERNAME"),
			os.Getenv("MYSQLPASSWORD"), os.Getenv("MYSQLDBNAME")))
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}

//GetUserById returns a user with specified id
func (db *Database) GetUserById(id int) (user models.User, err error) {
	return user, db.First(&user, id).Error
}

//RegisterUser adds a new user to database
func (db *Database) RegisterUser(user models.User) (models.User, error) {
	hashPassword(&user.Pass)
	user.LoggedIn = true
	return user, db.Create(&user).Error
}

// SignInUser validates user credentials and return his data if they are propper
func (db *Database) SignInUser(email, pass string) (user models.User, err error) {
	result := db.Table("users").Where(&models.User{Email: email})
	err = result.First(&user).Error
	if err != nil {
		return user, err
	}
	if !checkPassword(user.Pass, pass) {
		return user, errors.New("Invalid password")
	}
	user.Pass = ""
	err = result.Update("loggedin", 1).Error
	if err != nil {
		return user, err
	}

	return user, result.Find(&user).Error
}

// GetUserGroups returns a slice of Groups of which user is a member
func (db *Database) GetUserGroups(id int) (groups []models.Group, err error) {
	return groups, db.Table("members").Select("*").
		Joins("join members on group.id = group_id").
		Joins("join users on user.id = user_id").
		Where("user_id=?", id).Scan(&groups).Error
}

func hashPassword(s *string) error {
	if s == nil {
		return errors.New("Reference provided for hashing password is nil")
	}
	sBytes := []byte(*s)
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	*s = string(hashedBytes)
	return nil
}

func checkPassword(existingHash, incomingPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}
