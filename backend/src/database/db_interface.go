package database

import (
	"errors"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

type DBlayer interface {
	GetUserById(id int) (models.User, error)
	RegisterUser(models.User) (models.User, error)
	GetUserGroups(id uint) ([]models.Group, error)

	SignInUser(name string, pass string) (models.User, error)
	SignOutUser(id uint) error

	GetGroupMessages(id uint, offset uint) ([]models.Message, error)
	GetGroupMembership(id_group, id_user uint) (models.Member, error)

	CreateGroup(id uint, name, desc string) (models.Group, error)
	AddUserToGroup(name string, id_group uint, id_user uint) error
	DeleteUserFromGroup(id_member, id_group, id_user uint) error
	// GrantPriv(id_group, id_user int, priv models.Priv)
}

var ErrINVALIDPASSWORD = errors.New("invalid password")
var ErrNoPrivilages = errors.New("insufficient privilages")
var ErrInternal = errors.New("internal transaction error")
