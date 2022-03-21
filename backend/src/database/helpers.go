package database

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	Group   uint64    `json:"group"`
	Member  uint64    `json:"member"`
	Message string    `json:"text"`
	Nick    string    `json:"nick"`
	When    time.Time `json:"created"`
}

const MESSAGE_LIMIT = 4

func hashPassword(s string) (string, error) {
	if s == "" {
		return "", errors.New("Reference provided for hashing password is nil")
	}
	sBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	s = string(hashedBytes)
	return s, nil
}

func checkPassword(existingHash, incomingPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}
