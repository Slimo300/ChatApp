package repo

import (
	"errors"
	"time"
)

const TOKEN_VALID = "1"
const TOKEN_BLACKLISTED = "2"

type TokenRepository interface {
	SaveToken(token string, expiration time.Duration) error
	IsTokenValid(userID, tokenID string) (bool, error)
	InvalidateToken(userID, tokenID string) error
	InvalidateTokens(userID, tokenID string) error
}

var TokenBlacklistedError = errors.New("Token Blacklisted")
var TooManyTokensFoundError = errors.New("Too many tokens")
var TokenNotFoundError = errors.New("Token not found")
