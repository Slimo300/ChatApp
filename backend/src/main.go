package main

import (
	"log"

	"github.com/Slimo300/ChatApp/backend/src/database/orm"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/Slimo300/ChatApp/backend/src/routes"
	"github.com/Slimo300/ChatApp/backend/src/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	db, err := orm.SetupDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.Setup()
	server := handlers.NewServer(db, &storage)
	routes.Setup(engine, server)
	go server.RunHub()

	engine.Run(":8080")

}
