package handlers

import (
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetUserGroups(c *gin.Context) {
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "couldn't create group"})
		return
	}

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

func (s *Server) AddUserToGroup(c *gin.Context) {
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "no database connection"})
		return
	}

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
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "problems with server, try later"})
		return
	}

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
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	_, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "not authenticated"})
		return
	}

	load := struct {
		Group  int `json:"group"`
		Offset int `json:"offset"`
	}{}

	err = c.ShouldBindJSON(&load)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	messages, err := s.DB.GetGroupMessages(uint(load.Group), uint(load.Offset))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(messages) == 0 {
		c.JSON(http.StatusOK, load)
		return
	}

	c.JSON(http.StatusOK, messages)
}
