package handlers

import (
	"net/http"
	"strconv"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/gin-gonic/gin"
)

// handler for obtaining user membership for specified group
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

// handler for setting member privilages in his group
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

// handler for post method adding user to group
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

// handler for removing user from group doesn't delete membership just sets it to deleted, so that
// his messages would be still available for group
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
