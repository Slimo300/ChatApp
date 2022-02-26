package handlers

import (
	"github.com/Slimo300/ChatApp/backend/src/database"
)

type Server struct {
	DB     database.DBlayer
	secret string
	domain string
}

func NewServer() (*Server, error) {
	db, err := database.Setup()
	if err != nil {
		return nil, err
	}
	return &Server{DB: db, secret: "wołowina", domain: "localhost"}, nil
}
