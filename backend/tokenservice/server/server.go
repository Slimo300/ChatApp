package server

import (
	"crypto/rsa"
	"log"
	"time"

	"github.com/Slimo300/ChatApp/backend/tokensservice/repo"
)

type Service interface {
	NewPairFromUserID(userID string) (string, string, error)
	NewPairFromRefreshToken(refreshToken string) (string, string, error)
}

type tokenService struct {
	repo                  repo.TokenRepository
	refreshTokenSecret    string
	accessTokenPrivateKey rsa.PrivateKey
	accessTokenPublicKey  rsa.PublicKey
	refreshTokenDuration  time.Duration
	accessTokenDuration   time.Duration
}

func (t *tokenService) NewPairFromUserID(userID string) (refresh string, access string, err error) {
	refreshTokenData, err := t.generateRefreshToken(userID)
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}
	refresh = refreshTokenData.Token

	access, err = t.generateAccessToken(userID)
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	return
}
