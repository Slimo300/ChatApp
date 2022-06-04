package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
)

func TestGetUserGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := SetupTestServer()

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
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
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

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodGet, "/api/group/get", s.GetUserGroups)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			if tC.returnVal {
				respBody := []models.Group{}
				json.NewDecoder(response.Body).Decode(&respBody)

				if !reflect.DeepEqual(respBody, tC.expectedResponse) {
					t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
				}
			}
		})
	}
}

func TestDeleteGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := SetupTestServerWithHub()

	testCases := []struct {
		desc               string
		userID             uint
		groupID            string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		// user is not a creator of the group so he can't delete it
		{
			desc:               "deletegroupnosuccess",
			userID:             3,
			groupID:            "1",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "couldn't delete group"},
		},
		// user hasn't specified a group in a query
		{
			desc:               "deletegroupnotspecified",
			userID:             1,
			groupID:            "sa",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "group not specified"},
		},
		// creator deletes a group
		{
			desc:               "deletegroupsuccess",
			userID:             1,
			groupID:            "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.userID))
			if err != nil {
				t.Error("error when creating signed token")
			}
			req, _ := http.NewRequest(http.MethodDelete, "/api/group/"+tC.groupID, nil)
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodDelete, "/api/group/:groupID", s.DeleteGroup)
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

func TestCreateGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := SetupTestServer()

	testCases := []struct {
		desc               string
		ID                 uint
		data               map[string]interface{}
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		// no name provided in request body
		{
			desc:               "creategroupnoname",
			ID:                 3,
			data:               map[string]interface{}{"name": "", "desc": "ng1"},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "bad name"},
		},
		// no description provided in request body
		{
			desc:               "creategroupnodesc",
			ID:                 3,
			data:               map[string]interface{}{"name": "ng1", "desc": ""},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "bad description"},
		},
		// creator deletes a group
		{
			desc:               "creategroupsuccess",
			ID:                 3,
			data:               map[string]interface{}{"name": "ng1", "desc": "ng1"},
			returnVal:          true,
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   models.Group{ID: 2, Name: "ng1", Desc: "ng1"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.ID))
			if err != nil {
				t.Error("error when creating signed token")
			}
			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("POST", "/api/group/create", bytes.NewBuffer(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(s.MustAuth())
			engine.Handle(http.MethodPost, "/api/group/create", s.CreateGroup)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var respBody interface{}
			if tC.returnVal {
				member := models.Group{}
				json.NewDecoder(response.Body).Decode(&member)
				member.Created = time.Time{}
				respBody = member
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
