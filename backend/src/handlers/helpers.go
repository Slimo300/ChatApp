package handlers

import (
	"errors"
	"regexp"
)

var ErrNoDatabase = errors.New("No database connection")
var InvalidToken = errors.New("Invalid token")

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
