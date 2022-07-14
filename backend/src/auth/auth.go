package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/Slimo300/ChatApp/backend/tokenservice/pb"
)

type gRPCTokenAuthClient struct {
	client pb.TokenServiceClient
	pubKey rsa.PublicKey
}

func getPublicKey(pubKeyFile string) (*rsa.PublicKey, error) {
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	return pubKey, nil
}

func NewGRPCTokenAuthClient() *gRPCTokenAuthClient {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Println("Couldn't connect to grpc server: ", err.Error())
		return nil
	}
	client := pb.NewTokenServiceClient(conn)

	pubKey, err := getPublicKey(os.Getenv("PUB_KEY_FILE"))

	return &gRPCTokenAuthClient{
		client: client,
		pubKey: *pubKey,
	}
}

func (grpc *gRPCTokenAuthClient) NewPairFromUserID(userID uuid.UUID) (*pb.TokenPair, error) {
	response, err := grpc.client.NewPairFromUserID(context.Background(), &pb.UserID{ID: userID.String()})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (grpc *gRPCTokenAuthClient) NewPairFromRefresh(refresh string) (*pb.TokenPair, error) {
	response, err := grpc.client.NewPairFromRefresh(context.Background(), &pb.RefreshToken{Token: refresh})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (grpc *gRPCTokenAuthClient) DeleteUserToken(refresh string) error {
	response, err := grpc.client.DeleteUserToken(context.Background(), &pb.RefreshToken{Token: refresh})
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}
