package handlers_test

import (
	"github.com/Slimo300/ChatApp/backend/src/database/mock"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/Slimo300/ChatApp/backend/src/storage"
)

func SetupTestServerWithHub() handlers.Server {
	mockDB := mock.NewMockDB()
	s := handlers.NewServer(mockDB, storage.MockStorage{})
	go s.MockHub()
	return *s
}

func SetupTestServer() handlers.Server {
	mockDB := mock.NewMockDB()
	s := handlers.NewServer(mockDB, storage.MockStorage{})
	return *s
}
