package handlers

import (
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateGroup(c *gin.Context) {
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "couldn't create group"})
		return
	}

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	var group models.Group
	err = c.ShouldBindJSON(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	group, err = s.DB.CreateGroup(uint(id), group.Name, group.Desc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}
