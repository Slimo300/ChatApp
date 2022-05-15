package handlers_test

import (
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

func TestGetGroupMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := mock.NewMockDB()
	s := handlers.NewServer(mock)

	testCases := []struct {
		desc               string
		userID             uint
		returnVal          bool
		group              string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getmessagessuccess",
			userID:             1,
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			group:              "1",
			expectedResponse: []communication.Message{{Group: 1, Member: 1, Message: "elo", Nick: "Mal", When: "2019-13-01 22:00:45"},
				{Group: 1, Member: 2, Message: "siema", Nick: "River", When: "2019-15-01 22:00:45"},
				{Group: 1, Member: 1, Message: "elo elo", Nick: "Mal", When: "2019-16-01 22:00:45"},
				{Group: 1, Member: 2, Message: "siema siema", Nick: "River", When: "2019-17-01 22:00:45"}},
		},
		{
			desc:               "getmessagesunauthorized",
			userID:             3,
			returnVal:          false,
			group:              "1",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "User cannot request from this group"},
		},
		{
			desc:               "getmessagesnogroup",
			userID:             1,
			returnVal:          false,
			group:              "0",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.userID))
			if err != nil {
				t.Error("error when creating signed token")
			}

			req, _ := http.NewRequest("GET", "/api/group/"+tC.group+"/messages?num=4&offset=0", nil)
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodGet, "/api/group/:groupID/messages", s.GetGroupMessages)
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
