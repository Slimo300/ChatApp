package repo

import "time"

type TokenRepository interface {
	SaveToken(userID, tokenID string, expiration time.Duration) error
	IsTokenValid(userID, tokenID string) (bool, error)
	InvalidateToken(userID, tokenID string) error
}
