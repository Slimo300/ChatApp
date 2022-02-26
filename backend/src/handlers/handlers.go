package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (s *Server) Register(c *gin.Context) {
	if s.DB == nil {
		return
	}
	user := models.User{
		UserName: c.Query("name"),
		Email:    c.Query("email"),
		Pass:     c.Query("password"),
	}

	user, err := s.DB.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) SignIn(c *gin.Context) {
	if s.DB == nil {
		return
	}

	email := c.Query("email")
	password := c.Query("password")

	user, err := s.DB.SignInUser(email, password)
	if err != nil {
		if err == database.ErrINVALIDPASSWORD {
			c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.SetCookie("jwt", tokenString, 3600, "/", s.domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *Server) SignOutUser(c *gin.Context) {
	if s.DB == nil {
		return
	}

	if err := s.DB.SignOutUser(c.Query("email")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}

	c.SetCookie("jwt", "", -1, "/", s.domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
