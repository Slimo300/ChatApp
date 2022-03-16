package handlers

import (
	"strconv"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/ws"
	"github.com/golang-jwt/jwt"
)

type Server struct {
	DB     database.DBlayer
	Hub    *ws.Hub
	secret string
	domain string
}

func NewServer(db database.DBlayer) *Server {
	return &Server{
		DB:     db,
		secret: "wołowina",
		domain: "localhost",
		Hub:    ws.NewHub(),
	}
}

func (s *Server) RunHub() {
	s.Hub.Run()
}

func (s *Server) CreateSignedToken(iss int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(iss),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
