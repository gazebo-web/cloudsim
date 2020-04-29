package queue

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestIntegration_GetAll(t *testing.T) {
	userService := users.NewUserServiceMock()

	adminUsername := "root"
	admin := fuel.User{
		Name:             tools.Sptr("Admin Root"),
		Username:         &adminUsername,
		Email:            tools.Sptr("root@admin.com"),
	}

	userService.GetUserFromUsernameMock = func(username string) (user *fuel.User, msg *ign.ErrMsg) {
		return &admin, nil
	}

	userService.IsSystemAdminMock = func(user string) bool {
		return user == adminUsername
	}

	queue := NewQueue()
	service := NewService(queue, userService)
	controller := NewController(service)

	queue.Enqueue("1")
	queue.Enqueue("2")
	queue.Enqueue("3")

	router := mux.NewRouter()

	router.HandleFunc("/queue", func(writer http.ResponseWriter, request *http.Request) {
		result, err := controller.GetAll(&admin, writer, request)
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

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var response []string
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Len(t, response, 3)
}



func TestIntegration_MoveToFront(t *testing.T) {
	userService := users.NewUserServiceMock()

	adminUsername := "root"
	admin := fuel.User{
		Name:             tools.Sptr("Admin Root"),
		Username:         &adminUsername,
		Email:            tools.Sptr("root@admin.com"),
	}

	userService.GetUserFromUsernameMock = func(username string) (user *fuel.User, msg *ign.ErrMsg) {
		return &admin, nil
	}

	userService.IsSystemAdminMock = func(user string) bool {
		return user == adminUsername
	}

	queue := NewQueue()
	service := NewService(queue, userService)
	controller := NewController(service)

	router := mux.NewRouter()

	router.HandleFunc("/queue/{groupID}/front", func(writer http.ResponseWriter, request *http.Request) {
		result, err := controller.MoveToFront(&admin, writer, request)
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

	queue.Enqueue("1")
	queue.Enqueue("2")
	queue.Enqueue("3")

	req, err := http.NewRequest(http.MethodPatch, "/queue/3/front", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var response string
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, recorder.Code)

	items, _ := queue.Get(nil, nil)

	casted, ok := items[0].(string)
	assert.True(t, ok)
	assert.Equal(t, response, casted)
}



func TestIntegration_MoveToBack(t *testing.T) {
	userService := users.NewUserServiceMock()

	adminUsername := "root"
	admin := fuel.User{
		Name:             tools.Sptr("Admin Root"),
		Username:         &adminUsername,
		Email:            tools.Sptr("root@admin.com"),
	}

	userService.GetUserFromUsernameMock = func(username string) (user *fuel.User, msg *ign.ErrMsg) {
		return &admin, nil
	}

	userService.IsSystemAdminMock = func(user string) bool {
		return user == adminUsername
	}

	queue := NewQueue()
	service := NewService(queue, userService)
	controller := NewController(service)

	router := mux.NewRouter()

	router.HandleFunc("/queue/{groupID}/back", func(writer http.ResponseWriter, request *http.Request) {
		result, err := controller.MoveToBack(&admin, writer, request)
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

	queue.Enqueue("1")
	queue.Enqueue("2")
	queue.Enqueue("3")

	req, err := http.NewRequest(http.MethodPatch, "/queue/1/back", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var response string
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, recorder.Code)

	items, _ := queue.Get(nil, nil)

	casted, ok := items[2].(string)
	assert.True(t, ok)
	assert.Equal(t, response, casted)
}



func TestIntegration_Count(t *testing.T) {
	userService := users.NewUserServiceMock()

	adminUsername := "root"
	admin := fuel.User{
		Name:             tools.Sptr("Admin Root"),
		Username:         &adminUsername,
		Email:            tools.Sptr("root@admin.com"),
	}

	userService.GetUserFromUsernameMock = func(username string) (user *fuel.User, msg *ign.ErrMsg) {
		return &admin, nil
	}

	userService.IsSystemAdminMock = func(user string) bool {
		return user == adminUsername
	}

	queue := NewQueue()
	service := NewService(queue, userService)
	controller := NewController(service)

	router := mux.NewRouter()

	router.HandleFunc("/queue/count", func(writer http.ResponseWriter, request *http.Request) {
		result, err := controller.Count(&admin, writer, request)
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

	queue.Enqueue("1")
	queue.Enqueue("2")
	queue.Enqueue("3")

	req, err := http.NewRequest(http.MethodPatch, "/queue/count", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var response int
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, recorder.Code)

	count := queue.Count()
	assert.Equal(t, count, response)
}



func TestIntegration_Remove(t *testing.T) {
	userService := users.NewUserServiceMock()

	adminUsername := "root"
	admin := fuel.User{
		Name:             tools.Sptr("Admin Root"),
		Username:         &adminUsername,
		Email:            tools.Sptr("root@admin.com"),
	}

	userService.GetUserFromUsernameMock = func(username string) (user *fuel.User, msg *ign.ErrMsg) {
		return &admin, nil
	}

	userService.IsSystemAdminMock = func(user string) bool {
		return user == adminUsername
	}

	queue := NewQueue()
	service := NewService(queue, userService)
	controller := NewController(service)

	router := mux.NewRouter()

	router.HandleFunc("/queue/{groupID}", func(writer http.ResponseWriter, request *http.Request) {
		result, err := controller.Remove(&admin, writer, request)
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

	queue.Enqueue("1")
	queue.Enqueue("2")
	queue.Enqueue("3")

	req, err := http.NewRequest(http.MethodDelete, "/queue/1", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var response string
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, recorder.Code)

	items, _ := queue.Get(nil, nil)
	assert.Len(t, items, 2)
	assert.NotContains(t, items, response)
}
