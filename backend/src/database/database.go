package database

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

// Setup creates Database object and initializes connection between MySQL database
func Setup() (*Database, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@/%s?parseTime=true", os.Getenv("MYSQLUSERNAME"),
		os.Getenv("MYSQLPASSWORD"), os.Getenv("MYSQLDBNAME"))), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{}, models.Group{}, models.Member{}, models.Message{})

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

// GetUserGroups returns a slice of Groups of which user is a member
func (db *Database) GetUserGroups(id uint) (groups []models.Group, err error) {
	return groups, db.Table("`groups`").Select("`groups`.*").
		Joins("join `members` on `groups`.id = `members`.group_id").
		Joins("join `users` on `users`.id = `members`.user_id").
		Where("user_id=?", id).Scan(&groups).Error
}

func (db *Database) CreateGroup(id uint, name, desc string) (models.Group, error) {
	group := models.Group{Name: name, Desc: desc, Created: time.Now()}
	transactionFlag := false

	db.Transaction(func(tx *gorm.DB) error {
		creation := tx.Create(&group)
		if creation.Error != nil {
			return creation.Error
		}
		member := models.Member{UserID: id, GroupID: group.ID, Adding: true, Deleting: true, Setting: true, Creator: true}
		m_create := tx.Create(&member)
		if m_create.Error != nil {
			return m_create.Error
		}
		transactionFlag = true

		return nil
	})

	if transactionFlag == false {
		return models.Group{}, ErrInternal
	}

	return group, nil
}

func (db *Database) AddUserToGroup(username string, id_group uint, id_user uint) error {

	var member models.Member
	db.Table("`members`").Select("members.*").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Joins("inner join `groups` on `groups`.id = `members`.group_id").
		Where("`users`.id = ?", id_user).
		Where("`groups`.id = ?", id_group).Scan(&member)

	if !member.Adding {
		return ErrNoPrivilages
	}
	var user models.User
	selection := db.Where(&models.User{UserName: username}).First(&user)
	if selection.Error != nil {
		return selection.Error
	}
	member = models.Member{UserID: user.ID, Nick: user.UserName, GroupID: id_group, Adding: false, Deleting: false, Setting: false, Creator: false}
	creation := db.Create(&member)
	if creation.Error != nil {
		return creation.Error
	}

	return nil
}

func (db *Database) DeleteUserFromGroup(id_member, id_group, id_user uint) error {

	var member models.Member
	db.Table("`members`").Select("`members`.*").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Joins("inner join `groups` on `groups`.id = `members`.group_id").
		Where("`users`.id = ?", id_user).
		Where("`groups`.id = ?", id_group).Scan(&member)

	if !member.Deleting {
		return ErrNoPrivilages
	}

	var del_member models.Member
	selection := db.Where(&models.Member{ID: id_member}).First(&del_member)
	if selection.Error != nil {
		return selection.Error
	}

	deletion := db.Delete(del_member)
	if deletion.Error != nil {
		return deletion.Error
	}
	return nil
}

func (db *Database) GetGroupMessages(id uint, offset uint) ([]models.Message, error) {
	var messages []models.Message
	selection := db.Joins("Member", db.Where(&models.Member{GroupID: id})).Offset(int(offset)*15).Limit(15).Find(&messages, "`Member`.`group_id` = ?", id)
	if selection.Error != nil {
		return messages, selection.Error
	}

	return messages, nil
}

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
