package handlers_test

import (
	"bytes"
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
	s := handlers.NewServer(mock, nil)

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
	s := handlers.NewServer(mock, nil)

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

func TestDeleteGroup(t *testing.T) {
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
		// user is not a creator of the group so he can't delete it
		{
			desc:               "deletegroupnosuccess",
			ID:                 3,
			data:               map[string]interface{}{"group": 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Couldn't delete group"},
		},
		// user is dumb and hasn't specified a group in a query
		{
			desc:               "deletegroupnoquery",
			ID:                 1,
			data:               map[string]interface{}{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "group not specified"},
		},
		// creator deletes a group
		{
			desc:               "deletegroupsuccess",
			ID:                 1,
			data:               map[string]interface{}{"group": 1},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			jwt, err := s.CreateSignedToken(int(tC.ID))
			if err != nil {
				t.Error("error when creating signed token")
			}
			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("POST", "/api/group/delete", bytes.NewBuffer(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/group/delete", s.DeleteGroup)
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
	mock := database.NewMockDB()
	s := handlers.NewServer(mock, nil)

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
			expectedStatusCode: http.StatusOK,
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

func TestAddUser(t *testing.T) {
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
			desc:               "addusersuccess",
			ID:                 1,
			data:               map[string]interface{}{"username": "John", "group": 1},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
		// no name provided in request body
		{
			desc:               "addusernoname",
			ID:                 1,
			data:               map[string]interface{}{"group": 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "row not found"},
		},
		// no group provided in request body
		{
			desc:               "addusernogroup",
			ID:                 1,
			data:               map[string]interface{}{"username": "John"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "insufficient privilages"},
		},
		// user has no privilages to add to group
		{
			desc:               "addusernopriv",
			ID:                 2,
			data:               map[string]interface{}{"username": "John", "group": 1},
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

func TestDeleteUser(t *testing.T) {
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
			desc:               "deleteusersuccess",
			ID:                 1,
			data:               map[string]interface{}{"member": 2, "group": 1},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
		// no member provided in request body
		{
			desc:               "deletenomember",
			ID:                 1,
			data:               map[string]interface{}{"group": 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "row not found"},
		},
		// no member provided in request body
		{
			desc:               "deletenogroup",
			ID:                 1,
			data:               map[string]interface{}{"member": 2},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "insufficient privilages"},
		},
		// issuer has no right to add
		{
			desc:               "deletenpriv",
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
			req, _ := http.NewRequest("POST", "/api/group/remove", bytes.NewBuffer(requestBody))
			req.AddCookie(&http.Cookie{Name: "jwt", Value: jwt, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/group/remove", s.DeleteUserFromGroup)
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
