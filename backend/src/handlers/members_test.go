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
		// no member provided in request body
		{
			desc:               "grantprivnomember",
			userID:             1,
			memberID:           "100",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "row not found"},
		},
		// issuer has no right to add
		{
			desc:               "grantprivmemberdeleted",
			userID:             1,
			memberID:           "4",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "member deleted"},
		},
		{
			desc:               "grantprivnopriv",
			userID:             2,
			memberID:           "2",
			data:               map[string]interface{}{"adding": true, "deleting": true},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "insufficient privilages"},
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
			expectedResponse:   gin.H{"err": "member not specified"},
		},
		{
			desc:               "deletenopriv",
			userID:             2,
			memberID:           "2",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "insufficient privilages"},
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
