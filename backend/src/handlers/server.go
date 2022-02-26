package handlers

import (
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

type Server struct {
	DB database.DBlayer
}

func NewServer() (*Server, error) {
	db, err := database.Setup()
	if err != nil {
		return nil, err
	}
	return &Server{DB: db}, nil
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

	c.JSON(http.StatusOK, user)
}

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
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) SignOutUser(c *gin.Context) {
	if s.DB == nil {
		return
	}

	if err := s.DB.SignOutUser(c.Query("email")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
