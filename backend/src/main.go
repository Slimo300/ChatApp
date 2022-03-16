package main

import (
	"log"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/Slimo300/ChatApp/backend/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	// db := database.NewMockDB()
	db, err := database.Setup()
	if err != nil {
		log.Fatal(err)
	}
	server := handlers.NewServer(db)
	routes.Setup(engine, server)
	go server.RunHub()

	engine.Run(":8080")

}
