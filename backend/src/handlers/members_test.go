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
	"github.com/Slimo300/ChatApp/backend/src/database/mock"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/gin-gonic/gin"
)

func TestGrantPriv(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockDB := mock.NewMockDB()
	s := handlers.NewServer(mockDB, nil)

	testCases := []struct {
		desc               string
		userID             uint
		memberID           string
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "grantprivsuccess",
			userID:             1,
			memberID:           "2",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
		{
			desc:               "grantprivbadrequest",
			userID:             1,
			memberID:           "2",
			data:               map[string]interface{}{"adding": true, "deleting": true},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "bad request, all 3 fields must be present"},
		},
		// no member provided in request body
		{
			desc:               "grantprivnomember",
			userID:             1,
			memberID:           "100",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource not found"},
		},
		// issuer has no right to add
		{
			desc:               "grantprivmemberdeleted",
			userID:             1,
			memberID:           "4",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource not found"},
		},
		{
			desc:               "grantprivnorights",
			userID:             2,
			memberID:           "2",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": true},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to put"},
		},
		{
			desc:               "grantprivcreator",
			userID:             1,
			memberID:           "1",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "creator can't be modified"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.userID))
			if err != nil {
				t.Error("error when creating signed token")
			}

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest(http.MethodPut, "/api/member/"+tC.memberID, bytes.NewBuffer(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodPut, "/api/member/:memberID", s.GrantPriv)
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
	mockDB := mock.NewMockDB()

	channel := make(chan *communication.Action)
	go func() {
		for {
			<-channel
		}
	}()
	defer close(channel)
	s := handlers.NewServer(mockDB, channel)

	testCases := []struct {
		desc               string
		userID             uint
		data               map[string]interface{}
		memberID           string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "deleteusersuccess",
			userID:             1,
			memberID:           "2",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
		{
			desc:               "deletebadurl",
			userID:             1,
			memberID:           "0",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "member's id incorrect"},
		},
		{
			desc:               "deletenopriv",
			userID:             2,
			memberID:           "2",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to delete"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.userID))
			if err != nil {
				t.Error("error when creating signed token")
			}
			req, _ := http.NewRequest(http.MethodDelete, "/api/member/"+tC.memberID, nil)
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodDelete, "/api/member/:memberID", s.DeleteUserFromGroup)
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
