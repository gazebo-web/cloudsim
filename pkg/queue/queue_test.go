package queue

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueueIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

type IntegrationTestSuite struct {
	suite.Suite
	userService *users.ServiceMock
	adminUsername string
	admin fuel.User
	queue IQueue
	service IService
	controller IController
	router *mux.Router
	recorder *httptest.ResponseRecorder
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.userService = users.NewUserServiceMock()
	suite.adminUsername = "root"
	suite.admin = fuel.User{
		Name:             tools.Sptr("Admin Root"),
		Username:         &suite.adminUsername,
		Email:            tools.Sptr("root@admin.com"),
	}
	suite.queue = NewQueue()
	suite.service = NewService(suite.queue, suite.userService)
	suite.router = mux.NewRouter()
	suite.recorder = httptest.NewRecorder()
}


func (suite *IntegrationTestSuite) GetAll(t *testing.T) {
	suite.userService.On("GetUserFromUsername", suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", suite.admin.Username).Return(true)

	suite.queue.Enqueue("1")
	suite.queue.Enqueue("2")
	suite.queue.Enqueue("3")

	suite.router.HandleFunc("/queue", func(writer http.ResponseWriter, request *http.Request) {
		result, err := suite.controller.GetAll(&suite.admin, writer, request)
		if err != nil {
			body, _ := json.Marshal(err)
			writer.WriteHeader(err.StatusCode)
			writer.Write(body)
			return
		}
		body, _ := json.Marshal(result)
		writer.Write(body)
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	req, err := http.NewRequest(http.MethodGet, "/queue", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response []string
	assert.NoError(t, json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, suite.recorder.Code)
	assert.Len(t, response, 3)
}



func (suite *IntegrationTestSuite) MoveToFront(t *testing.T) {
	suite.userService.On("GetUserFromUsername", suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", suite.admin.Username).Return(true)

	suite.queue.Enqueue("1")
	suite.queue.Enqueue("2")
	suite.queue.Enqueue("3")

	suite.router.HandleFunc("/queue/{groupID}/front", func(writer http.ResponseWriter, request *http.Request) {
		result, err := suite.controller.MoveToFront(&suite.admin, writer, request)
		if err != nil {
			body, _ := json.Marshal(err)
			writer.WriteHeader(err.StatusCode)
			writer.Write(body)
			return
		}
		body, _ := json.Marshal(result)
		writer.Write(body)
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodPatch)

	req, err := http.NewRequest(http.MethodPatch, "/queue/3/front", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response string
	assert.NoError(t, json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, suite.recorder.Code)

	items, _ := suite.queue.Get(nil, nil)

	casted, ok := items[0].(string)
	assert.True(t, ok)
	assert.Equal(t, response, casted)
}



func (suite *IntegrationTestSuite) MoveToBack(t *testing.T) {
	suite.userService.On("GetUserFromUsername", suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", suite.admin.Username).Return(true)

	suite.queue.Enqueue("1")
	suite.queue.Enqueue("2")
	suite.queue.Enqueue("3")

	suite.router.HandleFunc("/queue/{groupID}/back", func(writer http.ResponseWriter, request *http.Request) {
		result, err := suite.controller.MoveToBack(&suite.admin, writer, request)
		if err != nil {
			body, _ := json.Marshal(err)
			writer.WriteHeader(err.StatusCode)
			writer.Write(body)
			return
		}
		body, _ := json.Marshal(result)
		writer.Write(body)
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodPatch)

	req, err := http.NewRequest(http.MethodPatch, "/queue/1/back", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response string
	assert.NoError(t, json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, suite.recorder.Code)

	items, _ := suite.queue.Get(nil, nil)

	casted, ok := items[2].(string)
	assert.True(t, ok)
	assert.Equal(t, response, casted)
}



func (suite *IntegrationTestSuite) Count(t *testing.T) {
	suite.userService.On("GetUserFromUsername", suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", suite.admin.Username).Return(true)

	suite.queue.Enqueue("1")
	suite.queue.Enqueue("2")
	suite.queue.Enqueue("3")

	suite.router.HandleFunc("/queue/count", func(writer http.ResponseWriter, request *http.Request) {
		result, err := suite.controller.Count(&suite.admin, writer, request)
		if err != nil {
			body, _ := json.Marshal(err)
			writer.WriteHeader(err.StatusCode)
			writer.Write(body)
			return
		}
		body, _ := json.Marshal(result)
		writer.Write(body)
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodPatch)

	req, err := http.NewRequest(http.MethodPatch, "/queue/count", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response int
	assert.NoError(t, json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, suite.recorder.Code)

	count := suite.queue.Count()
	assert.Equal(t, count, response)
}



func (suite *IntegrationTestSuite) Remove(t *testing.T) {
	suite.userService.On("GetUserFromUsername", suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", suite.admin.Username).Return(true)

	suite.queue.Enqueue("1")
	suite.queue.Enqueue("2")
	suite.queue.Enqueue("3")

	suite.router.HandleFunc("/queue/{groupID}", func(writer http.ResponseWriter, request *http.Request) {
		result, err := suite.controller.Remove(&suite.admin, writer, request)
		if err != nil {
			body, _ := json.Marshal(err)
			writer.WriteHeader(err.StatusCode)
			writer.Write(body)
			return
		}
		body, _ := json.Marshal(result)
		writer.Write(body)
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodDelete)

	req, err := http.NewRequest(http.MethodDelete, "/queue/1", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response string
	assert.NoError(t, json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, suite.recorder.Code)

	items, _ := suite.queue.Get(nil, nil)
	assert.Len(t, items, 2)
	assert.NotContains(t, items, response)
}
