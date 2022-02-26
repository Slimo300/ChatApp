package handlers

import (
	"github.com/Slimo300/ChatApp/backend/src/database"
)

type Server struct {
	DB     database.DBlayer
	secret string
	domain string
}

func NewServer(db database.DBlayer) (*Server, error) {
	return &Server{DB: db, secret: "wołowina", domain: "localhost"}, nil
}
