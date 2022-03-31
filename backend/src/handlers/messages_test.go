package handlers_test

import (
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

func TestGetGroupMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s := handlers.NewServer(mock, nil)

	testCases := []struct {
		desc               string
		data               uint
		returnVal          bool
		query              string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getmessagessuccess",
			data:               1,
			returnVal:          true,
			query:              "?group=1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: []communication.Message{{Group: 1, Member: 1, Message: "elo", Nick: "Mal", When: "2019-13-01 22:00:45"},
				{Group: 1, Member: 2, Message: "siema", Nick: "River", When: "2019-15-01 22:00:45"},
				{Group: 1, Member: 1, Message: "elo elo", Nick: "Mal", When: "2019-16-01 22:00:45"},
				{Group: 1, Member: 2, Message: "siema siema", Nick: "River", When: "2019-17-01 22:00:45"}},
		},
		{
			desc:               "getmessagesunauthorized",
			data:               3,
			returnVal:          false,
			query:              "?group=1",
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   gin.H{"err": "User cannot request from this group"},
		},
		{
			desc:               "getmessagesnogroup",
			data:               1,
			returnVal:          false,
			query:              "",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"message": "Select a group"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.data))
			if err != nil {
				t.Error("error when creating signed token")
			}

			req, _ := http.NewRequest("GET", "/api/group/messages"+tC.query, nil)
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodGet, "/api/group/messages", s.GetGroupMessages)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var respBody interface{}
			if tC.returnVal {
				groups := []communication.Message{}
				json.NewDecoder(response.Body).Decode(&groups)
				respBody = groups
			} else {
				var msg gin.H
				json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			json.NewDecoder(response.Body).Decode(&respBody)
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}
