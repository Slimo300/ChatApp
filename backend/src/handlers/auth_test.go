package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/auth"
	"github.com/Slimo300/ChatApp/backend/src/email"
	"github.com/Slimo300/ChatApp/backend/tokenservice/pb"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServer()
	mockEmailService := email.NewMockEmailService()
	s.EmailService = mockEmailService
	testCases := []struct {
		desc               string
		data               map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
		prepare            func(emailService *mock.Mock)
	}{
		{
			desc:               "registersuccess",
			data:               map[string]string{"username": "johnny", "email": "johnny@net.pl", "password": "password"},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"message": "success"},
			prepare: func(m *mock.Mock) {
				m.On("SendVerificationEmail", mock.Anything).Return(nil)
			},
		},
		{
			desc:               "registeremailtaken",
			data:               map[string]string{"username": "johnny12", "email": "johnny@net.pl", "password": "password"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "email already in database"},
			prepare:            func(emailService *mock.Mock) {},
		},
		{
			desc:               "registerinvalidpass",
			data:               map[string]string{"username": "johnny", "email": "johnny@net.pl", "password": ""},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid password"},
			prepare:            func(emailService *mock.Mock) {},
		},
		{
			desc:               "registerinvalidemail",
			data:               map[string]string{"username": "johnny", "email": "johnny@net.pl2", "password": "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email"},
			prepare:            func(emailService *mock.Mock) {},
		},
		{
			desc:               "registerinvalidusername",
			data:               map[string]string{"username": "j", "email": "johnny@net.pl", "password": "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid username"},
			prepare:            func(emailService *mock.Mock) {},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			tC.prepare(&mockEmailService.Mock)

			requestBody, _ := json.Marshal(tC.data)

			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/register", s.RegisterUser)
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
	mockTokenClient := auth.NewMockAuthClient()
	mockTokenClient.On("NewPairFromUserID", mock.Anything).Return(&pb.TokenPair{
		AccessToken:  "validAccessToken",
		RefreshToken: "validRefreshToken",
		Error:        "",
	}, nil)
	s := setupTestServerWithAuthClient(mockTokenClient)
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
			expectedResponse:   gin.H{"accessToken": "validAccessToken"},
		},
		{
			desc:               "loginnosuchemail",
			data:               map[string]string{"email": "mal.zein@email.co1m", "password": "test"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "wrong email or password"},
		},
		{
			desc:               "logininvalidpass",
			data:               map[string]string{"email": "mal.zein@email.com", "password": "t2est"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "wrong email or password"},
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
	mockAuthClient := auth.NewMockAuthClient()
	mockAuthClient.On("DeleteUserToken", mock.Anything).Return(nil)
	s := setupTestServerWithAuthClient(mockAuthClient)

	testCases := []struct {
		desc               string
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "logoutsuccess",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
		{
			desc:               "logoutnouser",
			id:                 "2f8fd072-29d4-470a-9359-b3b0e056bf65",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "No user with id: 2f8fd072-29d4-470a-9359-b3b0e056bf65"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/api/signout", nil)
			req.AddCookie(&http.Cookie{Name: "refreshToken", Value: "validRefreshToken", Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})
			w := httptest.NewRecorder()

			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.id)
			})

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

func TestRefresh(t *testing.T) {

	gin.SetMode(gin.TestMode)
	mockTokenClient := auth.NewMockAuthClient()
	s := setupTestServerWithAuthClient(mockTokenClient)

	testCases := []struct {
		desc               string
		withCookie         bool
		prepare            func(m *mock.Mock)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "refreshNoCookie",
			withCookie:         false,
			prepare:            func(m *mock.Mock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "No token provided"},
		},
		{
			desc:       "refreshTokenBlacklisted",
			withCookie: true,
			prepare: func(m *mock.Mock) {
				m.On("NewPairFromRefresh", mock.Anything).Return(&pb.TokenPair{
					AccessToken:  "",
					RefreshToken: "",
					Error:        "Token Blacklisted",
				}, nil).Once()
			},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Token Blacklisted"},
		},
		{
			desc:       "refreshOK",
			withCookie: true,
			prepare: func(m *mock.Mock) {
				m.On("NewPairFromRefresh", mock.Anything).Return(&pb.TokenPair{
					AccessToken:  "validAccessToken",
					RefreshToken: "validRefreshToken",
					Error:        "",
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"accessToken": "validAccessToken"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.prepare(&mockTokenClient.Mock)

			req := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
			if tC.withCookie {
				req.AddCookie(&http.Cookie{Name: "refreshToken", Value: "validRefreshToken", Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})
			}
			w := httptest.NewRecorder()

			_, engine := gin.CreateTestContext(w)

			engine.Handle(http.MethodPost, "/api/refresh", s.RefreshToken)
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
