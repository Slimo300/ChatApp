package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/handlers"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func TestGetUserGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s := handlers.NewServer(mock)

	date, _ := time.Parse("2006-01-02T15:04:05Z", "2019-01-13T08:47:44Z")

	testCases := []struct {
		desc               string
		data               uint
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getgroupssuccess",
			data:               1,
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []models.Group{{ID: 1, Name: "New Group", Desc: "totally new group", Created: date}},
		},
		{
			desc:               "getgroupsnone",
			data:               3,
			returnVal:          false,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "You don't have any group"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.data))
			if err != nil {
				t.Error("error when creating signed token")
			}

			req, _ := http.NewRequest("GET", "/api/group/get", nil)
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodGet, "/api/group/get", s.GetUserGroups)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var respBody interface{}
			if tC.returnVal {
				groups := []models.Group{}
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

func TestGetGroupMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s := handlers.NewServer(mock)

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
			expectedResponse: []database.Message{{Group: 1, Member: 1, Message: "elo", Nick: "Mal", When: "2019-13-01 22:00:45"},
				{Group: 1, Member: 2, Message: "siema", Nick: "River", When: "2019-15-01 22:00:45"},
				{Group: 1, Member: 1, Message: "elo elo", Nick: "Mal", When: "2019-16-01 22:00:45"},
				{Group: 1, Member: 2, Message: "siema siema", Nick: "River", When: "2019-17-01 22:00:45"}},
		},
		{
			desc:               "getmessagesunauthorized",
			data:               3,
			returnVal:          false,
			query:              "?group=1",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "User cannot request from this group"},
		},
		{
			desc:               "getmessagesunauthorized",
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
				groups := []database.Message{}
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

func TestGetMembership(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := database.NewMockDB()
	s := handlers.NewServer(mock)

	testCases := []struct {
		desc               string
		ID                 uint
		returnVal          bool
		query              string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getmembersuccess",
			ID:                 1,
			returnVal:          true,
			query:              "?group=1",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.Member{ID: 1, GroupID: 1, UserID: 1, Nick: "Mal", Adding: true, Deleting: true, Setting: true, Creator: true, Deleted: false},
		},
		{
			desc:               "getmembernosuch",
			ID:                 3,
			returnVal:          false,
			query:              "?group=1",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "Err no record"},
		},
		{
			desc:               "getmembernoquery",
			ID:                 3,
			returnVal:          false,
			query:              "",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"message": "Select a group"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.ID))
			if err != nil {
				t.Error("error when creating signed token")
			}

			req, _ := http.NewRequest("GET", "/api/group/membership"+tC.query, nil)
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodGet, "/api/group/membership", s.GetGroupMembership)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var respBody interface{}
			if tC.returnVal {
				member := models.Member{}
				json.NewDecoder(response.Body).Decode(&member)
				respBody = member
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
