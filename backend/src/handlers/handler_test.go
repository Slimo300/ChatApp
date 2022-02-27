package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s, err := handlers.NewServer(mock)
	if err != nil {
		t.Errorf("Couldn't create a server")
	}
	testCases := []struct {
		desc               string
		err                bool
		user               models.User
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "registersuccess",
			err:                false,
			user:               models.User{UserName: "johnny", Email: "johnny@net.pl", Pass: "password"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.User{UserName: "johnny", Email: "johnny@net.pl"},
		},
		{
			desc:               "registeremailtaken",
			err:                true,
			user:               models.User{UserName: "johnny", Email: "johnny@net.pl", Pass: "password"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   gin.H{"err": "couldn't register user"},
		},
		{
			desc:               "registerinvalidpass",
			err:                true,
			user:               models.User{UserName: "johnny", Email: "johnny@net.pl", Pass: ""},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid password"},
		},
		{
			desc:               "registerinvalidemail",
			err:                true,
			user:               models.User{UserName: "johnny", Email: "johnny@net.pl2", Pass: "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email"},
		},
		{
			desc:               "registerinvalidusername",
			err:                true,
			user:               models.User{UserName: "j", Email: "johnny@net.pl", Pass: "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid username"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			url := fmt.Sprintf("/api/register?name=%s&email=%s&password=%s", tC.user.UserName, tC.user.Email, tC.user.Pass)
			req := httptest.NewRequest(http.MethodPost, url, nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/register", s.Register)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var respBody interface{}
			if tC.err {
				var errmsg gin.H
				json.NewDecoder(response.Body).Decode(&errmsg)
				respBody = errmsg
			} else {
				var user models.User
				json.NewDecoder(response.Body).Decode(&user)
				respBody = models.User{UserName: user.UserName, Email: user.Email}
			}
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestSignIn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s, err := handlers.NewServer(mock)
	if err != nil {
		t.Errorf("Couldn't create a server")
	}
	testCases := []struct {
		desc               string
		err                bool
		user               models.User
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "loginsuccess",
			err:                false,
			user:               models.User{Email: "mal.zein@email.com", Pass: "test"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.User{Email: "mal.zein@email.com"},
		},
		{
			desc:               "logininvalidemail",
			err:                true,
			user:               models.User{Email: "mal.zein@email.com1", Pass: "test"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email"},
		},
		{
			desc:               "logininvalidpass",
			err:                true,
			user:               models.User{Email: "mal.zein@email.com", Pass: "passeord"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid password"},
		},
		{
			desc:               "loginnosuchuser",
			err:                true,
			user:               models.User{Email: "johnny@net.pl", Pass: "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			url := fmt.Sprintf("/api/login?&email=%s&password=%s", tC.user.Email, tC.user.Pass)
			req := httptest.NewRequest(http.MethodPost, url, nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/login", s.SignIn)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var respBody interface{}
			if tC.err {
				var errmsg gin.H
				json.NewDecoder(response.Body).Decode(&errmsg)
				respBody = errmsg
			} else {
				var user models.User
				json.NewDecoder(response.Body).Decode(&user)
				respBody = models.User{Email: user.Email}
			}
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}
