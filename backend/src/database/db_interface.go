package database

import (
	"errors"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/google/uuid"
)

type DBlayer interface {
	IsEmailInDatabase(email string) bool
	IsUsernameInDatabase(username string) bool

	GetUserById(uid uuid.UUID) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)

	RegisterUser(models.User) (models.User, error)
	SignInUser(id uuid.UUID) error
	SignOutUser(id uuid.UUID) error

	SetPassword(userID uuid.UUID, password string) error
	GetProfilePictureURL(userID uuid.UUID) (string, error)
	SetProfilePicture(userID uuid.UUID, newURI string) error
	DeleteProfilePicture(userID uuid.UUID) error

	GetUserGroups(id uuid.UUID) ([]models.Group, error)

	GetMemberByID(memberID uuid.UUID) (models.Member, error)
	GetUserGroupMember(userID, groupID uuid.UUID) (models.Member, error)

	GetGroupMessages(grouID uuid.UUID, offset, num int) ([]models.Message, error)
	AddMessage(memberID uuid.UUID, text string, when time.Time) error

	CreateGroup(ID uuid.UUID, name, desc string) (models.Group, error)
	DeleteUserFromGroup(memberID uuid.UUID) (models.Member, error)
	GrantPriv(memberID uuid.UUID, adding, deleting, setting bool) error

	DeleteGroup(groupID uuid.UUID) (models.Group, error)

	GetGroupProfilePicture(groupID uuid.UUID) (string, error)
	SetGroupProfilePicture(groupID uuid.UUID, newURI string) error
	DeleteGroupProfilePicture(groupID uuid.UUID) error

	GetUserInvites(userID uuid.UUID) ([]models.Invite, error)
	AddInvite(issID, targetID, groupID uuid.UUID) (models.Invite, error)

	DeclineInvite(inviteID uuid.UUID) error
	AcceptInvite(invite models.Invite) (models.Group, error)

	IsUserInGroup(userID, groupID uuid.UUID) bool
	IsUserInvited(userID, groupID uuid.UUID) bool

	GetInviteByID(inviteID uuid.UUID) (models.Invite, error)
}

const INVITE_AWAITING = 1
const INVITE_ACCEPT = 2
const INVITE_DECLINE = 3

const TIME_FORMAT = "2006-02-01 15:04:05"

var ErrINVALIDPASSWORD = errors.New("invalid password")
var ErrNoPrivilages = errors.New("insufficient privilages")
var ErrInternal = errors.New("internal transaction error")
