package handlers

import (
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterUser(c *gin.Context) {
	load := struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Pass     string `json:"password"`
	}{}
	err := c.ShouldBindJSON(&load)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}
	if !isEmailValid(load.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid email"})
		return
	}
	if !isPasswordValid(load.Pass) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid password"})
		return
	}
	if len(load.UserName) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid username"})
		return
	}
	if s.DB.IsUsernameInDatabase(load.UserName) {
		c.JSON(http.StatusConflict, gin.H{"err": "username taken"})
		return
	}
	if s.DB.IsEmailInDatabase(load.Email) {
		c.JSON(http.StatusConflict, gin.H{"err": "email already in database"})
		return
	}
	load.Pass, err = hashPassword(load.Pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	user := models.User{Email: load.Email, UserName: load.UserName, Pass: load.Pass}
	user, err = s.DB.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SignIn method
func (s *Server) SignIn(c *gin.Context) {
	load := struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}{}
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	requestPassword := load.Pass
	user, err := s.DB.GetUserByEmail(load.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "wrong email or password"})
		return
	}
	if !checkPassword(user.Pass, requestPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "wrong email or password"})
		return
	}
	if err := s.DB.SignInUser(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	// tokenPair, err := s.TokenService.NewPairFromUserID(user.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	// 	return
	// }
	token, err := s.CreateSignedToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.SetCookie("refreshToken", token, 3600, "/", s.domain, false, true)
	// c.SetCookie("refreshToken", tokenPair.RefreshToken, 3600, "/", s.domain, false, true)
	// c.Header("Authentication", tokenPair.AccessToken)

	c.JSON(http.StatusOK, gin.H{"name": user.UserName})
}

///////////////////////////////////////////////////////////////////////////////////////////////
// SignOutUser method
func (s *Server) SignOutUser(c *gin.Context) {
	userID := c.GetString("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
	}

	if err := s.DB.SignOutUser(uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.SetCookie("jwt", "", -1, "/", s.domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

///////////////////////////////////////////////////////////////////////////////////////////////
// GetUserById method
func (s *Server) GetUser(c *gin.Context) {
	userID := c.GetString("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
	}

	user, err := s.DB.GetUserById(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no such user"})
		return
	}
	user.Pass = ""

	c.JSON(http.StatusOK, user)
}
