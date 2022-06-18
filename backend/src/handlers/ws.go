package handlers

import (
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/ws"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) ServeWebSocket(c *gin.Context) {

	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	groups, err := s.DB.GetUserGroups(userUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	var grInt []uuid.UUID
	for _, group := range groups {
		grInt = append(grInt, group.ID)
	}

	ws.ServeWebSocket(c.Writer, c.Request, s.Hub, grInt, userUID)

}
