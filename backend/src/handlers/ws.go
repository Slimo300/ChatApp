package handlers

import (
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/ws"
	"github.com/gin-gonic/gin"
)

func (s *Server) ServeWebSocket(c *gin.Context) {
	if s.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	id, err := checkTokenAndGetID(c, s)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "not authenticated"})
		return
	}

	groups, err := s.DB.GetUserGroups(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	var grInt []int64
	for _, g := range groups {
		grInt = append(grInt, int64(g.ID))
	}

	ws.ServeWebSocket(c.Writer, c.Request, s.Hub, grInt)

}
