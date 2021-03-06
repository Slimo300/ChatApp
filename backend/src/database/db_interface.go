package database

import (
	"errors"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

type DBlayer interface {
	IsEmailInDatabase(email string) bool
	IsUsernameInDatabase(username string) bool

	GetUserById(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)

	RegisterUser(models.User) (models.User, error)
	SignInUser(id uint) error
	SignOutUser(id uint) error

	SetPassword(userID uint, password string) error
	GetProfilePictureURL(userID uint) (string, error)
	SetProfilePicture(userID uint, newURI string) error
	DeleteProfilePicture(userID uint) error

	GetUserGroups(id uint) ([]models.Group, error)

	GetMemberByID(memberID uint) (models.Member, error)
	GetUserGroupMember(userID, groupID uint) (models.Member, error)

	GetGroupMessages(id_group uint, offset, num int) ([]models.Message, error)
	AddMessage(memberID uint, text string, when time.Time) error

	CreateGroup(id uint, name, desc string) (models.Group, error)
	DeleteUserFromGroup(id_member uint) (models.Member, error)
	GrantPriv(id_member uint, adding, deleting, setting bool) error

	DeleteGroup(id_group uint) (models.Group, error)

	GetGroupProfilePicture(groupID uint) (string, error)
	SetGroupProfilePicture(groupID uint, newURI string) error
	DeleteGroupProfilePicture(groupID uint) error

	GetUserInvites(userID uint) ([]models.Invite, error)
	AddInvite(issID, targetID, groupID uint) (models.Invite, error)

	DeclineInvite(inviteID uint) error
	AcceptInvite(invite models.Invite) (models.Group, error)

	IsUserInGroup(userID, groupID uint) bool
	IsUserInvited(userID, groupID uint) bool

	GetInviteByID(inviteID uint) (models.Invite, error)
}

const INVITE_AWAITING = 1
const INVITE_ACCEPT = 2
const INVITE_DECLINE = 3

const TIME_FORMAT = "2006-02-01 15:04:05"

var ErrINVALIDPASSWORD = errors.New("invalid password")
var ErrNoPrivilages = errors.New("insufficient privilages")
var ErrInternal = errors.New("internal transaction error")
