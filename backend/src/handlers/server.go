package handlers

import (
	"log"
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

type Server struct {
	DB database.DBlayer
}

func NewServer() *Server {
	db, err := database.Setup()
	if err != nil {
		log.Fatalln("Error when creating server: ", err.Error())
	}
	return &Server{DB: db}
}

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
	user, err = s.DB.SignInUser(user.Email, user.Pass)
	if err != nil {
		if err == database.ErrINVALIDPASSWORD {
			c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}

	c.JSON(http.StatusOK, user)
}
