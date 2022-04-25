package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/gin-gonic/gin"
)

func TestGrantPriv(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s := handlers.NewServer(mock, nil)

	testCases := []struct {
		desc               string
		ID                 uint
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "grantprivsuccess",
			ID:                 1,
			data:               map[string]interface{}{"member": 2, "adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
		// no member provided in request body
		{
			desc:               "grantprivnomember",
			ID:                 1,
			data:               map[string]interface{}{"member": 100, "adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "row not found"},
		},
		// issuer has no right to add
		{
			desc:               "grantprivmemberdeleted",
			ID:                 1,
			data:               map[string]interface{}{"member": 4, "adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "member deleted"},
		},
		{
			desc:               "grantprivnopriv",
			ID:                 1,
			data:               map[string]interface{}{"member": 2, "adding": true, "deleting": true},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "bad request, all 3 fields must be present"},
		},
		{
			desc:               "grantprivcreator",
			ID:                 1,
			data:               map[string]interface{}{"member": 1, "adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "creator can't be modified"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.ID))
			if err != nil {
				t.Error("error when creating signed token")
			}
			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("PUT", "/api/group/rights", bytes.NewBuffer(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPut, "/api/group/rights", s.GrantPriv)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			if !reflect.DeepEqual(msg, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", msg, tC.expectedResponse)
			}
		})
	}
}

func TestAddMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()

	// mocking channel
	channel := make(chan *communication.Action)
	go func() {
		for {
			<-channel
		}
	}()
	defer close(channel)

	s := handlers.NewServer(mock, channel)

	testCases := []struct {
		desc               string
		ID                 uint
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "addusersuccess",
			ID:                 1,
			data:               map[string]interface{}{"username": "John", "group": 1},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"message": "ok"},
		},
		// no name provided in request body
		{
			desc:               "addusernoname",
			ID:                 1,
			data:               map[string]interface{}{"group": 1},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "row not found"},
		},
		// no group provided in request body
		{
			desc:               "addusernogroup",
			ID:                 1,
			data:               map[string]interface{}{"username": "John"},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   gin.H{"err": "insufficient privilages"},
		},
		// user has no privilages to add to group
		{
			desc:               "addusernopriv",
			ID:                 2,
			data:               map[string]interface{}{"username": "John", "group": 1},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   gin.H{"err": "insufficient privilages"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.ID))
			if err != nil {
				t.Error("error when creating signed token")
			}
			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("POST", "/api/group/add", bytes.NewBuffer(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/group/add", s.AddUserToGroup)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			if !reflect.DeepEqual(msg, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", msg, tC.expectedResponse)
			}
		})
	}
}

func TestDeleteMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()

	channel := make(chan *communication.Action)
	go func() {
		for {
			<-channel
		}
	}()
	defer close(channel)
	s := handlers.NewServer(mock, channel)

	testCases := []struct {
		desc               string
		ID                 uint
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "deleteusersuccess",
			ID:                 1,
			data:               map[string]interface{}{"member": 2, "group": 1},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
		// no member provided in request body
		{
			desc:               "deletenomember",
			ID:                 1,
			data:               map[string]interface{}{"group": 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "row not found"},
		},
		// issuer has no sufficient privilages to delete
		{
			desc:               "deletenopriv",
			ID:                 2,
			data:               map[string]interface{}{"member": 2, "group": 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "insufficient privilages"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.ID))
			if err != nil {
				t.Error("error when creating signed token")
			}
			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("PUT", "/api/group/remove", bytes.NewBuffer(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPut, "/api/group/remove", s.DeleteUserFromGroup)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			if !reflect.DeepEqual(msg, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", msg, tC.expectedResponse)
			}
		})
	}
}
