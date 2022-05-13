package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetGroupMessages(c *gin.Context) {
	userID := c.GetInt("userID")

	groupID := c.Param("groupID")
	if groupID == "" || groupID == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Select a group"})
		return
	}
	group_int, err := strconv.Atoi(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	num := c.Query("num")
	var num_int int
	if num == "" || num == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Set num parameter for number of messages you request"})
		return
	}
	num_int, err = strconv.Atoi(num)
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

	messages, err := s.DB.GetGroupMessages(uint(userID), uint(group_int), offset_int, num_int)
	if err != nil {
		if err.Error() == "User cannot request from this group" {
			c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if len(messages) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "no messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
