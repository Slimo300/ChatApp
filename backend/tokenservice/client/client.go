package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/Slimo300/ChatApp/backend/lib/pb"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewTokenServiceClient(conn)

	// response, err := c.NewPairFromUserID(context.Background(), &pb.UserID{ID: uuid.New().String()})
	// if err != nil {
	// 	log.Fatalf("Error when calling SayHello: %s", err)
	// }
	response, err := c.NewPairFromRefresh(context.Background(), &pb.RefreshToken{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTc1Mzk3MjksImp0aSI6ImIzNDkwOTk5LTgzMWQtNDRmZC1hMmQwLThmN2IxYTE3ZTM0ZCIsImlhdCI6MTY1NzQ1MzMyOSwic3ViIjoiOWI4NWE4NWEtYmM4Ni00MDkyLWFlNGItM2ZhZmYzN2Y1MmE2In0.ZNAoeVZl63Td6ad8QZarmM8HRu_TrJhM6RZFMFzA1BQ"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %v", response)

}
