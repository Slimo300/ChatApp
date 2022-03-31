package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const INVITE_AWAITING = 0
const INVITE_ACCEPT = 1
const INVITE_DECLINE = 2

const MESSAGE_LIMIT = 4
const TIME_FORMAT = "2006-02-01 15:04:05"

var ErrINVALIDPASSWORD = errors.New("invalid password")
var ErrNoPrivilages = errors.New("insufficient privilages")
var ErrInternal = errors.New("internal transaction error")

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
