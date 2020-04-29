package queue

import (
	"encoding/json"
	"github.com/gorilla/mux"
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
	suite.controller = NewController(suite.service)
}


func (suite *IntegrationTestSuite) TestGetAll() {
	suite.userService.On("GetUserFromUsername", *suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", *suite.admin.Username).Return(true)

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
		suite.Errorf(err, "Error creating HTTP Request.")
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response []string
	suite.NoError(json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Len(response, 3)
}

func (suite *IntegrationTestSuite) TestMoveToFront() {
	suite.userService.On("GetUserFromUsername", *suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", *suite.admin.Username).Return(true)

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
		suite.Errorf(err, "Error creating HTTP Request.")
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response string
	suite.NoError(json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	suite.Equal(http.StatusOK, suite.recorder.Code)

	items, _ := suite.queue.Get(nil, nil)

	casted, ok := items[0].(string)
	suite.True(ok)
	suite.Equal(response, casted)
}



func (suite *IntegrationTestSuite) TestMoveToBack() {
	suite.userService.On("GetUserFromUsername", *suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", *suite.admin.Username).Return(true)

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
		suite.Errorf(err, "Error creating HTTP Request.")
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response string
	suite.NoError(json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	suite.Equal(http.StatusOK, suite.recorder.Code)

	items, _ := suite.queue.Get(nil, nil)

	casted, ok := items[2].(string)
	suite.True(ok)
	suite.Equal(response, casted)
}



func (suite *IntegrationTestSuite) TestCount() {
	suite.userService.On("GetUserFromUsername", *suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", *suite.admin.Username).Return(true)

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
		suite.Errorf(err, "Error creating HTTP Request.")
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response int
	suite.NoError(json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	suite.Equal(http.StatusOK, suite.recorder.Code)

	count := suite.queue.Count()
	suite.Equal(count, response)
}



func (suite *IntegrationTestSuite) TestRemove() {
	suite.userService.On("GetUserFromUsername", *suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", *suite.admin.Username).Return(true)

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
		suite.Errorf(err,"Error creating HTTP Request.")
	}

	suite.router.ServeHTTP(suite.recorder, req)

	var response string
	suite.NoError(json.Unmarshal(suite.recorder.Body.Bytes(), &response))
	suite.Equal(http.StatusOK, suite.recorder.Code)

	items, _ := suite.queue.Get(nil, nil)
	suite.Len(items, 2)
	suite.NotContains(items, response)
}
