package handlers

import (
	"net/http"
	"strconv"

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
		c.JSON(http.StatusOK, gin.H{"message": "You don't have any group"})
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
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
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

	err = c.ShouldBindJSON(&load)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = s.DB.AddUserToGroup(load.Username, uint(load.Group), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (s *Server) DeleteUserFromGroup(c *gin.Context) {

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	load := struct {
		Group  int `json:"group"`
		Member int `json:"member"`
	}{}

	if err = c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err = s.DB.DeleteUserFromGroup(uint(load.Member), uint(load.Group), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *Server) GetGroupMessages(c *gin.Context) {

	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "not authenticated"})
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
		c.JSON(http.StatusForbidden, gin.H{"err": "not authenticated"})
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
		c.JSON(http.StatusForbidden, gin.H{"err": "not authenticated"})
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
	if err := s.DB.DeleteGroup(uint(load.Group), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}
