package handlers_test

import (
	"github.com/Slimo300/ChatApp/backend/lib/database/mock"
	"github.com/Slimo300/ChatApp/backend/lib/handlers"
	"github.com/Slimo300/ChatApp/backend/lib/storage"
)

func SetupTestServerWithHub() handlers.Server {
	mockDB := mock.NewMockDB()
	s := handlers.NewServerWithMockHub(mockDB, storage.MockStorage{})
	go s.RunHub()
	return *s
}

func SetupTestServer() handlers.Server {
	mockDB := mock.NewMockDB()
	s := handlers.NewServer(mockDB, storage.MockStorage{})
	return *s
}
