package handlers

import (
	"net/http"
	"strconv"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetUserGroups(c *gin.Context) {
	userID := c.GetInt("userID")

	groups, err := s.DB.GetUserGroups(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(groups) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, groups)

}

func (s *Server) CreateGroup(c *gin.Context) {
	userID := c.GetInt("userID")

	var group models.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if group.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad name"})
		return
	}
	if group.Desc == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad description"})
		return
	}

	group, err = s.DB.CreateGroup(uint(userID), group.Name, group.Desc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.sendHubChan <- &communication.Action{Group: int(group.ID), User: int(userID), Action: "CREATE_GROUP"}

	c.JSON(http.StatusCreated, group)
}

func (s *Server) DeleteGroup(c *gin.Context) {
	userID := c.GetInt("userID")

	groupID := c.Param("groupID")
	groupIDint, err := strconv.Atoi(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "group not specified"})
		return
	}

	member, err := s.DB.GetUserGroupMember(uint(userID), uint(groupIDint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "couldn't delete group"})
		return
	}
	if !member.Creator {
		c.JSON(http.StatusForbidden, gin.H{"err": "couldn't delete group"})
		return
	}

	group, err := s.DB.DeleteGroup(uint(groupIDint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	s.sendHubChan <- &communication.Action{Group: int(group.ID), Action: "DELETE_GROUP"}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}
