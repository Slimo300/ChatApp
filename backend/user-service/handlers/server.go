package handlers

import (
	"github.com/Slimo300/ChatApp/backend/src/auth"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/storage"
	"github.com/Slimo300/ChatApp/backend/user-service/email"
)

type Server struct {
	DB           database.DBlayer
	Storage      storage.StorageLayer
	TokenService auth.TokenClient
	EmailService email.EmailService
	domain       string
	MaxBodyBytes int64
}
