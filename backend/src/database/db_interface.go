package database

import (
	"github.com/Slimo300/ChatApp/backend/src/models"
)

type DBlayer interface {
	GetUserById(id int) (models.User, error)
	RegisterUser(models.User) (models.User, error)
	GetUserGroups(id uint) ([]models.Group, error)

	SignInUser(name string, pass string) (models.User, error)
	SignOutUser(id uint) error

	GetGroupMessages(id uint, offset uint) ([]Message, error)
	GetGroupMembership(id_group, id_user uint) (models.Member, error)
	// AddMessage adds a message of type Message to database
	AddMessage(msg Message) (Message, error)

	// AddFriend takes "id" of an issuer and "username" of invited user
	// AddFriend(id int, username string) (models.Invite, error)
	// RespondInvite takes id of an invite ("id") and response of type int (1 - agree, 2 - decline)
	// RespondInvite(id_inv, response int) (models.Group, error)

	CreateGroup(id uint, name, desc string) (models.Group, error)
	AddUserToGroup(name string, id_group uint, id_user uint) error
	DeleteUserFromGroup(id_member, id_group, id_user uint) error
	// GrantPriv(id_group, id_user int, priv models.Priv)

	// DeleteGroup deletes a specified group
	DeleteGroup(id_group, id_user uint) error
}
