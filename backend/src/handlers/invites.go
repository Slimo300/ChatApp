package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetUserInvites(c *gin.Context) {

	id := c.Value("userID").(int)

	invites, err := s.DB.GetUserInvites(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if len(invites) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, invites)
}

func (s *Server) SendGroupInvite(c *gin.Context) {
	userID := c.GetInt("userID")

	load := struct {
		GroupID int    `json:"group"`
		Target  string `json:"target"`
	}{}

	// getting req body
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if load.GroupID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "group not specified"})
		return
	}
	if strings.TrimSpace(load.Target) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not specified"})
		return
	}

	issuerMember, err := s.DB.GetUserGroupMember(uint(userID), uint(load.GroupID))
	if err != nil || !issuerMember.Adding {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to add"})
		return
	}

	userToBeAdded, err := s.DB.GetUserByUsername(load.Target)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": fmt.Sprintf("no user with name: %s", load.Target)})
		return
	}

	if s.DB.IsUserInGroup(userToBeAdded.ID, uint(load.GroupID)) {
		c.JSON(http.StatusConflict, gin.H{"err": "user is already a member of group"})
		return
	}

	if s.DB.IsUserInvited(userToBeAdded.ID, uint(load.GroupID)) {
		c.JSON(http.StatusConflict, gin.H{"err": "user already invited"})
		return
	}

	invite := models.Invite{IssId: uint(userID), TargetID: userToBeAdded.ID, GroupID: uint(load.GroupID)}

	s.CommChan <- &communication.Action{Invite: invite}

	c.JSON(http.StatusCreated, gin.H{"message": "invite sent"})
}

func (s *Server) RespondGroupInvite(c *gin.Context) {
	userID := c.GetInt("userID")
	inviteID := c.Param("inviteID")
	inviteIDint, err := strconv.Atoi(inviteID)
	if err != nil || inviteIDint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid invite id"})
		return
	}

	load := struct {
		Answer *bool `json:"answer" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "answer not specified"})
		return
	}

	invite, err := s.DB.GetInviteByID(uint(inviteIDint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "resource not found"})
		return
	}

	if invite.TargetID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to respond"})
		return
	}

	if invite.Status != database.INVITE_AWAITING {
		c.JSON(http.StatusForbidden, gin.H{"err": "invite already answered"})
		return
	}

	if !*load.Answer {
		if err := s.DB.DeclineInvite(invite.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "invite declined"})
		return
	}

	group, err := s.DB.AcceptInvite(invite)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no such invite"})
		return
	}

	c.JSON(http.StatusOK, group)
}
