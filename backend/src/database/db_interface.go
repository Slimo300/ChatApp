package database

import (
	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

type DBlayer interface {
	GetUserById(id int) (models.User, error)
	RegisterUser(models.User) (models.User, error)
	GetUserGroups(id uint) ([]models.Group, error)

	SignInUser(name string, pass string) (models.User, error)
	SignOutUser(id uint) error

	GetGroupMessages(id_user, id_group uint, offset, num int) ([]communication.Message, error)
	AddMessage(msg communication.Message) (communication.Message, error)

	CreateGroup(id uint, name, desc string) (models.Group, error)
	AddUserToGroup(name string, id_group uint, id_user uint) (models.Member, error)
	DeleteUserFromGroup(id_member, id_user uint) (models.Member, error)
	GrantPriv(id_member, id_user uint, adding, deleting, setting bool) error

	DeleteGroup(id_group, id_user uint) (models.Group, error)

	SendGroupInvite(issId, groupID uint, target string) (models.Invite, error)
}
