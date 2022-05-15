package handlers

import (
	"net/http"
	"strconv"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetGroupMessages(c *gin.Context) {
	userID := c.GetInt("userID")

	groupID := c.Param("groupID")
	groupIDint, err := strconv.Atoi(groupID)
	if err != nil || groupIDint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}

	num := c.Query("num")
	numInt, err := strconv.Atoi(num)
	if err != nil || numInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "number of messages is not a valid number"})
		return
	}

	offset := c.Query("offset")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "offset is not a valid number"})
		return
	}

	if !s.DB.IsUserInGroup(uint(userID), uint(groupIDint)) {
		c.JSON(http.StatusForbidden, gin.H{"err": "User cannot request from this group"})
		return
	}

	messages, err := s.DB.GetGroupMessages(uint(groupIDint), offsetInt, numInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if len(messages) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	lightMessages := communication.ShortenMessages(messages)

	c.JSON(http.StatusOK, lightMessages)
}
