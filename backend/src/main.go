package main

import (
	"log"

	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/Slimo300/ChatApp/backend/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	routes.Setup(engine, server)

	engine.Run(":8080")

}
