package handlers

import (
	"net/http"
	"strings"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Server) SendGroupInvite(c *gin.Context) {
	// getting id from jwt
	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "not authenticated"})
		return
	}

	load := struct {
		Group  int    `json:"group"`
		Target string `json:"target"`
	}{}

	// getting req body
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if load.Group == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "group not specified"})
		return
	}
	if strings.TrimSpace(load.Target) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not specified"})
		return
	}

	invite, err := s.DB.SendGroupInvite(uint(id), uint(load.Group), load.Target)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
		return
	}

	// Send invite via channel
	s.CommChan <- &communication.Action{Invite: invite}

	c.JSON(http.StatusCreated, gin.H{"message": "invite sent"})
}
