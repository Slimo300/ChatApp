package routes

import (
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, server *handlers.Server) {
	engine.Use(CORSMiddleware())
	engine.POST("/api/login", server.SignIn)
	engine.POST("/api/register", server.Register)
	engine.POST("/api/signout", server.SignOutUser)
}
