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

	// getting group from request
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

	// getting number of messages to fetch from request
	num := c.Query("num")
	var num_int int
	if num == "" || num == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Set num parameter for number of messages you request"})
		return
	} else {
		num_int, err = strconv.Atoi(num)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		}
	}

	// getting offset from request
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

	// getting messages from database
	messages, err := s.DB.GetGroupMessages(uint(id), uint(group_int), offset_int, num_int)
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
