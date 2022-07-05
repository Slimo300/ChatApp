package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Slimo300/ChatApp/backend/tokensservice/pb"
	"github.com/Slimo300/ChatApp/backend/tokensservice/server"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := server.Server{}

	grpcServer := grpc.NewServer()

	pb.RegisterTokenServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
