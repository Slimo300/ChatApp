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

	db.AutoMigrate(&models.User{}, models.Group{}, models.Member{}, models.Priv{})

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

func (db *Database) CreateGroup(id uint, name, desc string) (models.Group, error) {
	group := models.Group{Name: name, Desc: desc, Created: time.Now()}

	db.Transaction(func(tx *gorm.DB) error {
		creation := tx.Create(&group)
		if creation.Error != nil {
			return creation.Error
		}
		member := models.Member{UserID: id, GroupID: group.ID, PrivID: 9}
		m_create := tx.Create(&member)
		if m_create.Error != nil {
			return m_create.Error
		}

		return nil
	})

	return group, nil
}

func (db *Database) AddUserToGroup(username string, id_group uint, id_user uint) error {

	var priv models.Priv
	db.Table("members").Select("priv.adding").Joins("inner join priv on priv.id = members.id_priv").
		Joins("inner join users on users.id = members.user_id").
		Joins("inner join groups on groups.id = members.group_id").
		Where("users.id = ?", id_user).
		Where("groups.id = ?", id_group).Scan(&priv)

	if !priv.Adding {
		return ErrNoPrivilages
	}
	var user models.User
	selection := db.Where(&models.User{UserName: username}).First(&user)
	if selection.Error != nil {
		return selection.Error
	}
	member := models.Member{UserID: user.ID, GroupID: id_group, PrivID: 1}
	creation := db.Create(&member)
	if creation.Error != nil {
		return creation.Error
	}

	return nil
}

func (db *Database) DeleteUserFromGroup(id_member, id_group, id_user uint) error {

	var priv models.Priv
	db.Table("members").Select("priv.adding").Joins("inner join priv on priv.id = members.id_priv").
		Joins("inner join users on users.id = members.user_id").
		Joins("inner join groups on groups.id = members.group_id").
		Where("users.id = ?", id_user).
		Where("groups.id = ?", id_group).Scan(&priv)

	if !priv.Deleting {
		return ErrNoPrivilages
	}

	var member models.Member
	selection := db.Where(&models.Member{ID: id_member}).First(&member)
	if selection.Error != nil {
		return selection.Error
	}

	deletion := db.Delete(&member)
	if deletion.Error != nil {
		return deletion.Error
	}
	return nil
}

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
