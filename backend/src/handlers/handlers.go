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
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "couldn't register user"})
		return
	}
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err0": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"err": "couldn't register user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SignIn method
func (s *Server) SignIn(c *gin.Context) {
	if s.DB == nil {
		return
	}
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

	user, err = s.DB.SignInUser(user.Email, user.Pass)
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

	c.JSON(http.StatusOK, gin.H{"name": user.UserName})
}

///////////////////////////////////////////////////////////////////////////////////////////////
// SignOutUser method
func (s *Server) SignOutUser(c *gin.Context) {
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	if err = s.DB.SignOutUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.SetCookie("jwt", "", -1, "/", s.domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

///////////////////////////////////////////////////////////////////////////////////////////////
// GetUserById method
func (s *Server) GetUserById(c *gin.Context) {
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}

	user, err := s.DB.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "no such user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
