package routes

import (
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, server *handlers.Server) {
	engine.Use(CORSMiddleware())
	engine.Use(server.CheckDatabase())
	engine.POST("/api/login", server.SignIn)
	engine.POST("/api/register", server.Register)
	engine.POST("/api/signout", server.SignOutUser)
	engine.GET("/api/user", server.GetUser)
	engine.POST("/api/group/create", server.CreateGroup)
	engine.POST("/api/group/add", server.AddUserToGroup)
	engine.POST("/api/group/remove", server.DeleteUserFromGroup)
	engine.POST("/api/group/delete", server.DeleteGroup)
	engine.GET("/api/group/get", server.GetUserGroups)
	engine.GET("/api/group/messages", server.GetGroupMessages)
	engine.GET("/api/group/membership", server.GetGroupMembership)
	engine.GET("/ws", server.ServeWebSocket)
}
