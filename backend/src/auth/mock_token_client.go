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

func NewMockAuthClient() MockAuthClient {
	return MockAuthClient{}
}

func (m MockAuthClient) NewPairFromUserID(userID uuid.UUID) (*pb.TokenPair, error) {
	ret := m.Called(userID)
	return ret.Get(0).(*pb.TokenPair), ret.Error(1)
}
func (m MockAuthClient) NewPairFromRefresh(refresh string) (*pb.TokenPair, error) {
	ret := m.Called(refresh)
	return ret.Get(0).(*pb.TokenPair), ret.Error(1)
}

func (m MockAuthClient) DeleteUserToken(refresh string) error {
	ret := m.Called(refresh)
	return ret.Error(0)
}

func (m MockAuthClient) GetPublicKey() *rsa.PublicKey {
	ret := m.Called()
	return ret.Get(0).(*rsa.PublicKey)
}
