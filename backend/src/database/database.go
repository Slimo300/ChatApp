package database

import (
	"errors"
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

	// Getting ids of groups in which user has membership
	var usergroups []uint
	if err := db.Table("`groups`").Select("`groups`.id").
		Joins("inner join `members` on `members`.group_id = `groups`.id").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Where("`members`.deleted = false").
		Where("`users`.id = ?", id).Scan(&usergroups).Error; err != nil {
		return groups, err
	}

	// Getting full groups data with its members
	if err := db.Where("id in (?)", usergroups).Preload("Members", "deleted is false").Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
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

func (db *Database) AddUserToGroup(username string, id_group uint, id_user uint) (models.Member, error) {

	var member models.Member
	if err := db.Table("`members`").Select("members.*").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Joins("inner join `groups` on `groups`.id = `members`.group_id").
		Where("`users`.id = ?", id_user).
		Where("`groups`.id = ?", id_group).Scan(&member).Error; err != nil {
		return models.Member{}, err
	}

	if !member.Adding {
		return models.Member{}, ErrNoPrivilages
	}

	var user models.User
	if err := db.Where(&models.User{UserName: username}).First(&user).Error; err != nil {
		return models.Member{}, err
	}

	var member2 models.Member
	if err := db.Where(&models.Member{UserID: user.ID, GroupID: id_group}).First(&member2).Error; err != nil && err != gorm.ErrRecordNotFound {
		return models.Member{}, err
	}
	// if member does not exist member.Deleted is false
	if member2.Deleted == true {
		if err := db.Model(member2).Update("deleted", false).Error; err != nil {
			return models.Member{}, err
		}
		return models.Member{}, nil
	}

	member = models.Member{UserID: user.ID, Nick: user.UserName, GroupID: id_group, Adding: false, Deleting: false, Setting: false, Creator: false}
	if err := db.Create(&member).Error; err != nil {
		return models.Member{}, err
	}

	return member, nil
}

func (db *Database) DeleteUserFromGroup(id_member, id_user uint) (models.Member, error) {

	var deleted_member models.Member
	if err := db.First(&deleted_member, id_member).Error; err != nil {
		return models.Member{}, err
	}

	var issuer_member models.Member
	if err := db.Table("`members`").Select("`members`.*").
		Joins("inner join `users` on `users`.id = `members`.user_id").
		Joins("inner join `groups` on `groups`.id = `members`.group_id").
		Where("`users`.id = ?", id_user).
		Where("`groups`.id = ?", deleted_member.GroupID).Scan(&issuer_member).Error; err != nil {
		return models.Member{}, err
	}

	if !issuer_member.Deleting {
		return models.Member{}, ErrNoPrivilages
	}

	if err := db.Model(&deleted_member).Update("deleted", true).Error; err != nil {
		return models.Member{}, err
	}

	return deleted_member, nil
}

func (db *Database) GetGroupMembership(id_group, id_user uint) (models.Member, error) {
	var membership models.Member
	selection := db.Where(&models.Member{UserID: id_user, GroupID: id_group}).First(&membership)
	if selection.Error != nil {
		return membership, selection.Error
	}

	return membership, nil
}

// Deletes a specified group if user is authorized to do so
func (db *Database) DeleteGroup(id_group, id_user uint) (models.Group, error) {
	// getting user membership to check his privilages
	var membership models.Member
	if err := db.Where(&models.Member{UserID: id_user, GroupID: id_group}).First(&membership).Error; err != nil {
		return models.Group{}, err
	}

	// checking whether user have privilages to delete a group
	if !membership.Creator {
		return models.Group{}, ErrNoPrivilages
	}

	// deleting memberships
	if err := db.Where(models.Member{GroupID: id_group}).Delete(&models.Member{}).Error; err != nil {
		return models.Group{}, err
	}
	// deleting specified group
	group := models.Group{ID: id_group}
	if err := db.Delete(&group).Error; err != nil {
		return models.Group{}, err
	}

	return group, nil
}

// Updates user rights to a group
func (db *Database) GrantPriv(id_mem, id uint, adding, deleting, setting bool) error {

	var member models.Member
	if err := db.First(&member, id_mem).Error; err != nil {
		return err
	}
	// can't modify deleted member nor a creator one
	if member.Deleted {
		return errors.New("member deleted")
	}
	if member.Creator {
		return errors.New("creator can't be modified")
	}

	// checking membership of an issuer
	var issuer models.Member
	if err := db.Where(models.Member{UserID: id, GroupID: member.GroupID}).First(&issuer).Error; err != nil {
		return err
	}
	if !issuer.Setting {
		return ErrNoPrivilages
	}

	if err := db.Model(member).Updates(models.Member{Adding: adding, Deleting: deleting, Setting: setting}).Error; err != nil {
		return err
	}

	return nil
}
