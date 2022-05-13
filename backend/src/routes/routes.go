package routes

import (
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, server *handlers.Server) {

	api := engine.Group("api")

	api.Use(CORSMiddleware())
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

	apiAuth.GET("/ws", server.ServeWebSocket)

	apiAuth.GET("/invites", server.GetUserInvites)
	apiAuth.POST("/invites", server.SendGroupInvite)
	apiAuth.PUT("/invites/:inviteID", server.RespondGroupInvite)
}
