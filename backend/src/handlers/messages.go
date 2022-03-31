package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
