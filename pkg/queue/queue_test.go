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
	})

	req, err := http.NewRequest(http.MethodGet, "/queue", nil)
	if err != nil {
		t.Errorf("Error creating HTTP Request. Error: [%v]", err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var response []Item
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestIntegration_Get(t *testing.T) {
}

func TestIntegration_Enqueue(t *testing.T) {
}

func TestIntegration_Dequeue(t *testing.T) {
}

func TestIntegration_EnqueueOrWait(t *testing.T) {
}

func TestIntegration_MoveToFront(t *testing.T) {

}

func TestIntegration_MoveToBack(t *testing.T) {
}

func TestIntegration_Count(t *testing.T) {
}

func TestIntegration_Remove(t *testing.T) {
}
