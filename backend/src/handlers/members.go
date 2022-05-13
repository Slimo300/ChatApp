package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) GrantPriv(c *gin.Context) {
	userID := c.GetInt("userID")
	memberID := c.Param("memberID")
	if memberID == "" || memberID == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "member not specified"})
		return
	}
	memberIDint, err := strconv.Atoi(memberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "member id not specified")
		return
	}

	var member models.Member

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad request, all 3 fields must be present"})
		return
	}

	if err := s.DB.GrantPriv(uint(memberIDint), uint(userID), member.Adding, member.Deleting, member.Setting); err != nil {
		if err.Error() == "creator can't be modified" {
			c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
			return
		}
		if errors.Is(err, database.ErrNoPrivilages) {
			c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (s *Server) DeleteUserFromGroup(c *gin.Context) {
	userID := c.GetInt("userID")
	memberID := c.Param("memberID")
	if memberID == "" || memberID == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "member not specified"})
		return
	}
	memberIDint, err := strconv.Atoi(memberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "member id not correct"})
		return
	}

	member, err := s.DB.DeleteUserFromGroup(uint(memberIDint), uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	s.CommChan <- &communication.Action{Action: "DELETE_MEMBER", Member: member}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
