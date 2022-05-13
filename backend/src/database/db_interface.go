package database

import (
	"errors"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

type DBlayer interface {
	IsEmailInDatabase(email string) bool
	IsUsernameInDatabase(username string) bool

	GetUserById(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)

	RegisterUser(models.User) (models.User, error)
	SignInUser(id uint) error
	SignOutUser(id uint) error

	GetUserGroups(id uint) ([]models.Group, error)

	GetGroupMessages(id_user, id_group uint, offset, num int) ([]communication.Message, error)
	AddMessage(msg communication.Message) (communication.Message, error)

	CreateGroup(id uint, name, desc string) (models.Group, error)
	AddUserToGroup(name string, id_group uint, id_user uint) (models.Member, error)
	DeleteUserFromGroup(id_member, id_user uint) (models.Member, error)
	GrantPriv(id_member, id_user uint, adding, deleting, setting bool) error

	DeleteGroup(id_group, id_user uint) (models.Group, error)

	GetUserInvites(userID uint) ([]models.Invite, error)
	SendGroupInvite(issId, groupID uint, target string) (models.Invite, error)
	RespondGroupInvite(userID, inviteID uint, response bool) (models.Group, error)
}

const INVITE_AWAITING = 0
const INVITE_ACCEPT = 1
const INVITE_DECLINE = 2

const TIME_FORMAT = "2006-02-01 15:04:05"

var ErrINVALIDPASSWORD = errors.New("invalid password")
var ErrNoPrivilages = errors.New("insufficient privilages")
var ErrInternal = errors.New("internal transaction error")
