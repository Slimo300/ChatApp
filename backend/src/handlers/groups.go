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
	id := c.GetInt("userID")

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
	id := c.GetInt("userID")

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

	group, err = s.DB.CreateGroup(uint(id), group.Name, group.Desc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.CommChan <- &communication.Action{Group: int(group.ID), User: int(id), Action: "CREATE_GROUP"}

	c.JSON(http.StatusCreated, group)
}

func (s *Server) DeleteGroup(c *gin.Context) {
	id := c.GetInt("userID")

	groupID := c.Param("groupID")
	if groupID == "0" || groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "group not specified"})
		return
	}
	groupIDint, err := strconv.Atoi(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "something went wrong"})
	}

	// telling database to delete group
	group, err := s.DB.DeleteGroup(uint(groupIDint), uint(id))
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

	s.CommChan <- &communication.Action{Group: int(group.ID), Action: "DELETE_GROUP"}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}
