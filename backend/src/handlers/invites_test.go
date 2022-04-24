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

func TestSendGroupInvite(t *testing.T) {
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
			expectedResponse:   gin.H{"err": "user not found"},
		},
		{
			desc:               "invitenoprivilages",
			id:                 2,
			data:               map[string]interface{}{"group": 1, "target": "Kel"},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "insufficient privilages"},
		},
		{
			desc:               "inviteuserismember",
			id:                 1,
			data:               map[string]interface{}{"group": 1, "target": "River"},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "user already in a group"},
		},
		{
			desc:               "invitealreadyindatabase",
			id:                 1,
			data:               map[string]interface{}{"group": 1, "target": "John"},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "invite already sent"},
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
			engine.Handle(http.MethodPost, "/api/invite", s.SendGroupInvite)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var respBody interface{}
			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)
			respBody = msg

			json.NewDecoder(response.Body).Decode(&respBody)
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}
