package handlers

import (
	"net/http"
	"strconv"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/gin-gonic/gin"
)

func (s *Server) GrantPriv(c *gin.Context) {
	userID := c.GetInt("userID")
	memberID := c.Param("memberID")
	memberIDint, err := strconv.Atoi(memberID)
	if err != nil || memberIDint <= 0 {
		c.JSON(http.StatusBadRequest, "member's id incorrect")
		return
	}

	load := struct {
		Adding   *bool `json:"adding" binding:"required"`
		Deleting *bool `json:"deleting" binding:"required"`
		Setting  *bool `json:"setting" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad request, all 3 fields must be present"})
		return
	}

	memberToBeChanged, err := s.DB.GetMemberByID(uint(memberIDint))
	if err != nil || memberToBeChanged.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"err": "resource not found"})
		return
	}

	issuerMember, err := s.DB.GetUserGroupMember(uint(userID), memberToBeChanged.GroupID)
	if err != nil || !issuerMember.Setting {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to put"})
		return
	}

	if memberToBeChanged.Creator {
		c.JSON(http.StatusForbidden, gin.H{"err": "creator can't be modified"})
	}

	if err := s.DB.GrantPriv(uint(memberIDint), *load.Adding, *load.Deleting, *load.Setting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (s *Server) DeleteUserFromGroup(c *gin.Context) {
	userID := c.GetInt("userID")
	memberID := c.Param("memberID")
	memberIDint, err := strconv.Atoi(memberID)
	if err != nil || memberIDint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "member's id incorrect"})
		return
	}

	memberToBeDeleted, err := s.DB.GetMemberByID(uint(memberIDint))
	if err != nil || memberToBeDeleted.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"err": "resource not found"})
		return
	}

	issuerMember, err := s.DB.GetUserGroupMember(uint(userID), memberToBeDeleted.GroupID)
	if err != nil || (!issuerMember.Deleting && issuerMember.ID != memberToBeDeleted.ID) {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to delete"})
		return
	}

	member, err := s.DB.DeleteUserFromGroup(uint(memberIDint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.actionChan <- &communication.Action{Action: "DELETE_MEMBER", Member: member}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
