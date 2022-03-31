package handlers

import (
	"net/http"
	"strconv"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetUserGroups(c *gin.Context) {

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	groups, err := s.DB.GetUserGroups(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(groups) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "You don't have any group"})
		return
	}

	c.JSON(http.StatusOK, groups)

}

func (s *Server) CreateGroup(c *gin.Context) {

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
	if group.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad name"})
		return
	}
	if group.Desc == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad description"})
		return
	}

	group, err = s.DB.CreateGroup(uint(id), group.Name, group.Desc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.CommChan <- &communication.Action{Group: int(group.ID), User: int(id), Action: "insert"}

	c.JSON(http.StatusCreated, group)
}

func (s *Server) AddUserToGroup(c *gin.Context) {

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	load := struct {
		Username string `json:"username"`
		Group    int    `json:"group"`
	}{}

	if err = c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	member, err := s.DB.AddUserToGroup(load.Username, uint(load.Group), uint(id))
	if err != nil {
		if err.Error() == "insufficient privilages" {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		if err.Error() == "row not found" {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.CommChan <- &communication.Action{Group: int(member.GroupID), User: int(member.UserID), Action: "insert"}

	c.JSON(http.StatusCreated, gin.H{"message": "ok"})
}

func (s *Server) DeleteUserFromGroup(c *gin.Context) {

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	load := struct {
		Member int `json:"member"`
	}{}

	if err = c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// TODO channel to hub
	_, err = s.DB.DeleteUserFromGroup(uint(load.Member), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (s *Server) GetGroupMessages(c *gin.Context) {

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "not authenticated"})
		return
	}

	group := c.Query("group")
	if group == "" || group == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Select a group"})
		return
	}
	group_int, err := strconv.Atoi(group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	offset := c.Query("offset")
	var offset_int int
	if offset == "" {
		offset_int = 0
	} else {
		offset_int, err = strconv.Atoi(offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
	}

	messages, err := s.DB.GetGroupMessages(uint(id), uint(group_int), uint(offset_int))
	if err != nil {
		if err.Error() == "User cannot request from this group" {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(messages) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (s *Server) GetGroupMembership(c *gin.Context) {

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "not authenticated"})
		return
	}

	group := c.Query("group")
	if group == "" || group == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Select a group"})
		return
	}

	group_int, err := strconv.Atoi(group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	member, err := s.DB.GetGroupMembership(uint(group_int), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

func (s *Server) DeleteGroup(c *gin.Context) {
	// getting id from jwt
	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "not authenticated"})
		return
	}

	load := struct {
		Group int `json:"group"`
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

	// telling database to delete group
	group, err := s.DB.DeleteGroup(uint(load.Group), uint(id))
	if err != nil {
		if err.Error() == "Couldn't delete group" {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
			return
		}
		if err == database.ErrNoPrivilages {
			c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.CommChan <- &communication.Action{User: 0, Group: int(group.ID), Action: "pop"}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}

func (s *Server) GrantPriv(c *gin.Context) {
	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	load := struct {
		Member   int   `json:"member" binding:"required"`
		Adding   *bool `json:"adding" binding:"required"`
		Deleting *bool `json:"deleting" binding:"required"`
		Setting  *bool `json:"setting" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad request, all 3 fields must be present"})
		return
	}

	if err := s.DB.GrantPriv(uint(load.Member), uint(id), *load.Adding, *load.Deleting, *load.Setting); err != nil {
		if err.Error() == "creator can't be modified" {
			c.JSON(http.StatusForbidden, gin.H{"err": "creator can't be modified"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
