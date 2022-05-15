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
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func TestSendGroupInvite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockDB := mock.NewMockDB()

	// mocking channel
	channel := make(chan *communication.Action)
	go func() {
		for {
			<-channel
		}
	}()
	defer close(channel)

	s := handlers.NewServer(mockDB)

	testCases := []struct {
		desc               string
		id                 int
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "invitesuccess",
			id:                 1,
			data:               map[string]interface{}{"group": 1, "target": "Kel"},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"message": "invite sent"},
		},
		{
			desc:               "invitenosuchuser",
			id:                 1,
			data:               map[string]interface{}{"group": 1, "target": "Raul"},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "no user with name: Raul"},
		},
		{
			desc:               "invitenorights",
			id:                 2,
			data:               map[string]interface{}{"group": 1, "target": "Kel"},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to add"},
		},
		{
			desc:               "inviteuserismember",
			id:                 1,
			data:               map[string]interface{}{"group": 1, "target": "River"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "user is already a member of group"},
		},
		{
			desc:               "invitealreadyindatabase",
			id:                 1,
			data:               map[string]interface{}{"group": 1, "target": "John"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "user already invited"},
		},
		{
			desc:               "invitenogroup",
			id:                 1,
			data:               map[string]interface{}{"target": "Kel"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "group not specified"},
		},
		{
			desc:               "invitenouser",
			id:                 1,
			data:               map[string]interface{}{"group": 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "user not specified"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(tC.id)
			if err != nil {
				t.Error("error when creating signed token")
			}
			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("POST", "/api/invite", bytes.NewReader(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodPost, "/api/invite", s.SendGroupInvite)
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

func TestGetUserInvites(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockDB := mock.NewMockDB()

	s := handlers.NewServer(mockDB)

	dateCreated, _ := time.Parse("2006-01-02T15:04:05Z", "2019-03-17T22:04:45Z")
	dateModified, _ := time.Parse("2006-01-02T15:04:05Z", "2019-03-17T22:04:45Z")

	testCases := []struct {
		desc               string
		id                 int
		expectedStatusCode int
		expectedResponse   []models.Invite
	}{
		{
			desc:               "getinvitessuccess",
			id:                 3,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []models.Invite{{ID: 1, IssId: 1, TargetID: 3, GroupID: 1, Status: 0, Created: dateCreated, Modified: dateModified}},
		},
		{
			desc:               "getinvitesnocontent",
			id:                 1,
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   []models.Invite{},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(tC.id)
			if err != nil {
				t.Error("error when creating signed token")
			}

			req, _ := http.NewRequest("GET", "/api/invites", nil)
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodGet, "/api/invites", s.GetUserInvites)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			respBody := []models.Invite{}
			json.NewDecoder(response.Body).Decode(&respBody)

			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestRespondGroupInvite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockDB := mock.NewMockDB()

	s := handlers.NewServer(mockDB)

	dateGroupCreated, _ := time.Parse("2006-01-02T15:04:05Z", "2019-01-13T08:47:44Z")

	testCases := []struct {
		desc               string
		userID             int
		answer             bool
		data               map[string]interface{}
		inviteID           string
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "respondInviteYes",
			userID:             3,
			data:               map[string]interface{}{"answer": true},
			inviteID:           "1",
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.Group{ID: 1, Name: "New Group", Desc: "totally new group", Created: dateGroupCreated},
		},
		{
			desc:               "respondInviteNo",
			userID:             3,
			data:               map[string]interface{}{"answer": false},
			inviteID:           "1",
			returnVal:          false,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "invite declined"},
		},
		{
			desc:               "respondInviteNotInDatabase",
			userID:             3,
			data:               map[string]interface{}{"answer": true},
			inviteID:           "2",
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource not found"},
		},
		{
			desc:               "respondInviteWrongUser",
			userID:             1,
			data:               map[string]interface{}{"answer": true},
			inviteID:           "1",
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to respond"},
		},
		{
			desc:               "respondInviteNoAnswer",
			userID:             1,
			inviteID:           "1",
			data:               map[string]interface{}{},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "answer not specified"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(tC.userID)
			if err != nil {
				t.Error("error when creating signed token")
			}

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("PUT", "/api/invite/"+tC.inviteID, bytes.NewReader(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodPut, "/api/invite/:inviteID", s.RespondGroupInvite)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var respBody interface{}
			if tC.returnVal {
				group := models.Group{}
				json.NewDecoder(response.Body).Decode(&group)
				respBody = group
			} else {
				var msg gin.H
				json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}
