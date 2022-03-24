package database

import (
	"fmt"
	"os"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
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
	// Updating logged field
	if err = result.Update("logged", 1).Error; err != nil {
		return user, err
	}
	// Updating activity field
	if err = result.Update("activity", time.Now()).Error; err != nil {
		return user, err
	}

	return user, nil
}

// SignOutUser updates user "logged" field
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

	var issuer models.User
	db.Where(models.User{ID: id}).First(&issuer)

	db.Transaction(func(tx *gorm.DB) error {
		creation := tx.Create(&group)
		if creation.Error != nil {
			return creation.Error
		}
		member := models.Member{UserID: id, GroupID: group.ID, Adding: true, Deleting: true, Setting: true, Creator: true, Nick: issuer.UserName}
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

	deletion := db.Delete(&del_member)
	if deletion.Error != nil {
		return deletion.Error
	}
	return nil
}

// GetGroupMembership returns Member model stating user and group relation and user rights in it
func (db *Database) GetGroupMembership(id_group, id_user uint) (models.Member, error) {
	var membership models.Member
	if err := db.Where(&models.Member{UserID: id_user, GroupID: id_group}).First(&membership).Error; err != nil {
		return membership, err
	}

	return membership, nil
}

// Deletes a specified group if user is authorized to do so
func (db *Database) DeleteGroup(id_group, id_user uint) error {
	// getting user membership to check his privilages
	var membership models.Member
	if err := db.Where(&models.Member{UserID: id_user, GroupID: id_group}).First(&membership).Error; err != nil {
		return err
	}

	// checking whether user have privilages to delete a group
	if !membership.Creator {
		return ErrNoPrivilages
	}

	// deleting specified group
	if err := db.Delete(&models.Group{ID: id_group}).Error; err != nil {
		return err
	}

	return nil
}
