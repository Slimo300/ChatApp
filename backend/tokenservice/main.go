package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/Slimo300/ChatApp/backend/lib/pb"
	"github.com/Slimo300/ChatApp/backend/lib/repo/redis"
	"github.com/Slimo300/ChatApp/backend/lib/server"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Go gRPC Beginners Tutorial!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	priv, err := ioutil.ReadFile(os.Getenv("PRIV_KEY_FILE"))
	if err != nil {
		log.Fatal("could not read private key pem file: %w", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		log.Fatal("could not parse private key: %w", err)
	}

	repo := redis.NewRedisTokenRepository("localhost", "6379", "")

	s := server.NewTokenService(repo, "wolowina", *privKey, 24*time.Hour, 20*time.Minute)

	grpcServer := grpc.NewServer()

	pb.RegisterTokenServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
