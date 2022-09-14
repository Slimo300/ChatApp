package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/auth"
	"github.com/Slimo300/ChatApp/backend/src/database/orm"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/Slimo300/ChatApp/backend/src/routes"
	"github.com/Slimo300/ChatApp/backend/src/storage"
	"github.com/Slimo300/ChatApp/backend/user-service/email"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	db, err := orm.Setup()
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.Setup()
	authService, err := auth.NewGRPCTokenAuthClient()
	if err != nil {
		panic("Couldn't connect to grpc auth server")
	}
	server := handlers.NewServer(db, &storage, authService)
	server.EmailService = email.NewMockEmailService()
	routes.Setup(engine, server)
	go server.RunHub()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

}
