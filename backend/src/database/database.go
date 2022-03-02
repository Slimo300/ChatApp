package database

import (
	"errors"
	"fmt"
	"os"
	"time"

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

	db.AutoMigrate(&models.User{})

	return &Database{DB: db}, nil
}

//GetUserById returns a user with specified id
func (db *Database) GetUserById(id int) (user models.User, err error) {
	return user, db.First(&user, id).Error
}

//RegisterUser adds a new user to database
func (db *Database) RegisterUser(user models.User) (models.User, error) {
	pass, err := hashPassword(user.Pass)
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
	if !checkPassword(user.Pass, pass) {
		return user, ErrINVALIDPASSWORD
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
	err := user.Update("logged", 0).Error
	if err != nil {
		return err
	}
	return nil
}

// // GetUserGroups returns a slice of Groups of which user is a member
// func (db *Database) GetUserGroups(id int) (groups []models.Group, err error) {
// 	return groups, db.Table("members").Select("*").
// 		Joins("join members on group.id = group_id").
// 		Joins("join users on user.id = user_id").
// 		Where("user_id=?", id).Scan(&groups).Error
// }

// func (db *Database) AddUserToGroup(id int) error {
// 	return nil
// }

// func (db *Database) CreateGroup(name, desc string) (models.Group, error) {
// 	return models.Group{ID: 1, Name: name, Desc: desc, Created: time.Now()}, nil
// }

// func (db *Database) DeleteUserFromGroup(id int) error {
// 	return nil
// }

// func (db *Database) GetGroupMessages(id int, since time.Time) ([]models.Message, error) {
// 	return nil, nil
// }

// func (db *Database) GrantPriv(granter models.User, receiver models.User, err error)

func hashPassword(s string) (string, error) {
	if s == "" {
		return "", errors.New("Reference provided for hashing password is nil")
	}
	sBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	s = string(hashedBytes)
	return s, nil
}

func checkPassword(existingHash, incomingPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}
