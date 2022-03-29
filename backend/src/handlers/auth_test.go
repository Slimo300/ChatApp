package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/gin-gonic/gin"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s := handlers.NewServer(mock, nil)
	testCases := []struct {
		desc               string
		data               map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "registersuccess",
			data:               map[string]string{"username": "johnny", "email": "johnny@net.pl", "password": "password"},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"message": "success"},
		},
		{
			desc:               "registeremailtaken",
			data:               map[string]string{"username": "johnny", "email": "johnny@net.pl", "password": "password"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "email taken"},
		},
		{
			desc:               "registerinvalidpass",
			data:               map[string]string{"username": "johnny", "email": "johnny@net.pl", "password": ""},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid password"},
		},
		{
			desc:               "registerinvalidemail",
			data:               map[string]string{"username": "johnny", "email": "johnny@net.pl2", "password": "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email"},
		},
		{
			desc:               "registerinvalidusername",
			data:               map[string]string{"username": "j", "email": "johnny@net.pl", "password": "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid username"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			requestBody, _ := json.Marshal(tC.data)

			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/register", s.Register)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var errmsg gin.H
			json.NewDecoder(response.Body).Decode(&errmsg)
			if !reflect.DeepEqual(errmsg, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", errmsg, tC.expectedResponse)
			}
		})
	}
}

func TestSignIn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s := handlers.NewServer(mock, nil)
	testCases := []struct {
		desc               string
		data               map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "loginsuccess",
			data:               map[string]string{"email": "mal.zein@email.com", "password": "test"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"name": "Mal"},
		},
		{
			desc:               "logininvalidemail",
			data:               map[string]string{"email": "mal.zein@email.co1m", "password": "test"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email"},
		},
		{
			desc:               "logininvalidpass",
			data:               map[string]string{"email": "mal.zein@email.com", "password": "t2est"},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   gin.H{"err": "invalid password"},
		},
		{
			desc:               "loginnosuchuser",
			data:               map[string]string{"email": "mal2.zein@email.com", "password": "test"},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "No email mal2.zein@email.com in database"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			requestBody, _ := json.Marshal(tC.data)
			req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(requestBody))
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/login", s.SignIn)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestSignOut(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := database.NewMockDB()
	s := handlers.NewServer(db, nil)

	testCases := []struct {
		desc               string
		id                 int
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "logoutsuccess",
			id:                 1,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
		{
			desc:               "logoutnouser",
			id:                 1000,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "No user with id: 1000"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			w := httptest.NewRecorder()

			jwt, err := s.CreateSignedToken(tC.id)
			if err != nil {
				t.Error("error when creating signed token")
			}

			_, engine := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodPost, "/api/signout", nil)
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			engine.Handle(http.MethodPost, "/api/signout", s.SignOutUser)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}
