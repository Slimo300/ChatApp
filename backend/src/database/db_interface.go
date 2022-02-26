package database

import (
	"errors"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

type DBlayer interface {
	GetUserById(id int) (models.User, error)
	RegisterUser(models.User) (models.User, error)
	// GetUserGroups(id int) ([]models.Group, error)

	SignInUser(name string, pass string) (models.User, error)
	SignOutUser(email string) error

	// GetGroupMessages(id int, since time.Time) ([]models.Message, error)

	// CreateGroup(name, desc string) (models.Group, error)
	// AddUserToGroup(id int) error
	// DeleteUserFromGroup(id int) error
	// GrantPriv(id_group, id_user int, priv models.Priv)
}

var ErrINVALIDPASSWORD = errors.New("Invalid password")
