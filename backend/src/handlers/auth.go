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

//////////////////////////////////////////////////////////////////////////////////////////////////
// Register method
func (s *Server) Register(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}
	if !isEmailValid(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid email"})
		return
	}
	if len(user.Pass) <= 6 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid password"})
		return
	}
	if len(user.UserName) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid username"})
		return
	}

	user, err = s.DB.RegisterUser(user)
	if err != nil {
		if err.Error() == "email taken" || err.Error() == "username taken" {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SignIn method
func (s *Server) SignIn(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if !isEmailValid(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid email"})
		return
	}
	email := user.Email
	user, err = s.DB.SignInUser(user.Email, user.Pass)
	if err != nil {
		if err == database.ErrINVALIDPASSWORD {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		if err.Error() == "No email "+email+" in database" {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{"name": user.UserName})
}

///////////////////////////////////////////////////////////////////////////////////////////////
// SignOutUser method
func (s *Server) SignOutUser(c *gin.Context) {
	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	if err = s.DB.SignOutUser(uint(id)); err != nil {
		if err.Error() == "No user with id: "+strconv.Itoa(id) {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.SetCookie("jwt", "", -1, "/", s.domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

///////////////////////////////////////////////////////////////////////////////////////////////
// GetUserById method
func (s *Server) GetUser(c *gin.Context) {
	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := s.DB.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no such user"})
		return
	}
	user.Pass = ""
	user.ID = 0

	c.JSON(http.StatusOK, user)
}
