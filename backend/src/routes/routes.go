package routes

import (
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, server *handlers.Server) {

	engine.Use(CORSMiddleware())

	api := engine.Group("/api")
	api.Use(server.CheckDatabase())

	api.POST("/register", server.RegisterUser)
	api.POST("/login", server.SignIn)

	apiAuth := api.Use(server.MustAuth())
	apiAuth.POST("/signout", server.SignOutUser)
	apiAuth.GET("/user", server.GetUser)

	apiAuth.GET("/group", server.GetUserGroups)
	apiAuth.POST("/group", server.CreateGroup)
	apiAuth.DELETE("/group/:groupID", server.DeleteGroup)

	apiAuth.DELETE("/member/:memberID", server.DeleteUserFromGroup)
	apiAuth.PUT("/member/:memberID", server.GrantPriv)

	apiAuth.GET("/group/:groupID/messages", server.GetGroupMessages)

	apiAuth.GET("/invites", server.GetUserInvites)
	apiAuth.POST("/invites", server.SendGroupInvite)
	apiAuth.PUT("/invites/:inviteID", server.RespondGroupInvite)

	ws := engine.Group("/ws")
	ws.Use(server.MustAuth())
	ws.GET("/", server.ServeWebSocket)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
