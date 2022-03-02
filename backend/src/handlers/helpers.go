package handlers

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func checkTokenAndGetID(c *gin.Context, s *Server) (int, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return 0, err
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secret), nil
		})
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(token.Claims.(*jwt.StandardClaims).Issuer)
	if err != nil {
		return 0, errors.New(("Bad token"))
	}

	return id, nil

}
