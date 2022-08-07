package auth

import (
	"crypto/rsa"

	"github.com/Slimo300/ChatApp/backend/tokenservice/pb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAuthClient struct {
	mock.Mock
}

func (MockAuthClient) NewPairFromUserID(userID uuid.UUID) (*pb.TokenPair, error)
func (MockAuthClient) NewPairFromRefresh(refresh string) (*pb.TokenPair, error)
func (MockAuthClient) DeleteUserToken(refresh string) error
func (MockAuthClient) GetPublicKey() *rsa.PublicKey
