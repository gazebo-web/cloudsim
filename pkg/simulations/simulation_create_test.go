package simulations

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSimulationIntegration(t *testing.T) {
	suite.Run(t, new(simulationTestSuite))
}

type simulationTestSuite struct {
	suite.Suite
	userService   *users.ServiceMock
	adminUsername string
	admin         fuel.User

	repository *RepositoryMock
	service    Service
	controller Controller

	router   *mux.Router
	recorder *httptest.ResponseRecorder

	updateSimulation *SimulationUpdate
	createSimulation *SimulationCreate
	simulation       *Simulation
}

func (suite *simulationTestSuite) SetupSuite() {

	suite.userService = users.NewServiceMock()

	suite.repository = NewRepositoryMock()
	suite.service = NewService(suite.repository)

	input := NewControllerInput{
		SimulationService: suite.service,
		UserService:       suite.userService,
		FormDecoder:       form.NewDecoder(),
		Validator:         validator.New(),
	}
	suite.controller = NewController(input)

	suite.router = mux.NewRouter()
	suite.recorder = httptest.NewRecorder()

	suite.adminUsername = "root"
	suite.admin = fuel.User{
		Name:     tools.Sptr("Admin Root"),
		Username: &suite.adminUsername,
		Email:    tools.Sptr("root@admin.com"),
	}
}

func (suite *simulationTestSuite) TestCreate() {

	suite.router.HandleFunc("/simulation", func(writer http.ResponseWriter, request *http.Request) {
		result, err := suite.controller.Start(&suite.admin, writer, request)
		if err != nil {
			body, _ := json.Marshal(err)
			writer.WriteHeader(err.StatusCode)
			writer.Write(body)
			return
		}
		body, _ := json.Marshal(result)
		writer.Write(body)
		writer.WriteHeader(http.StatusOK)
	}).Methods(http.MethodPost)

	createForm := map[string]string{
		"name":        "test-1234",
		"owner":       suite.adminUsername,
		"robot_name":  "X1",
		"robot_type":  "X1_SENSOR_CONFIG_1",
		"robot_image": "infrastructureascode/aws-cli:latest",
	}

	body := bytes.NewBuffer(nil)
	formWriter := multipart.NewWriter(body)

	for key, value := range createForm {
		suite.NoError(formWriter.WriteField(key, value))
	}
	suite.NoError(formWriter.Close())

	result := Simulation{
		ID:            1,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		DeletedAt:     nil,
		StoppedAt:     nil,
		ValidFor:      tools.Sptr("1m"),
		Owner:         nil,
		Creator:       &suite.adminUsername,
		Private:       nil,
		StopOnEnd:     nil,
		Name:          tools.Sptr("test-name"),
		Image:         tools.Sptr("image-test"),
		GroupID:       tools.Sptr("aaaa-bbbb-cccc-dddd-eeee"),
		ParentGroupID: nil,
		MultiSim:      0,
		Status:        StatusPending.ToIntPtr(),
		ErrorStatus:   nil,
		Platform:      tools.Sptr("cloudsim"),
		Application:   tools.Sptr("test"),
		Robots:        nil,
		Held:          false,
		Extra:         tools.Sptr("extra-field"),
		ExtraSelector: tools.Sptr("extra-selector"),
	}

	suite.userService.On("GetUserFromUsername", *suite.admin.Username).Return(suite.admin, nil)
	suite.userService.On("IsSystemAdmin", *suite.admin.Username).Return(true)
	suite.userService.On("AddResourcePermission", result.GroupID, per.Read).Return(true, nil)
	suite.userService.On("AddResourcePermission", result.GroupID, per.Write).Return(true, nil)
	suite.userService.On("VerifyOwner", *suite.admin.Username, suite.adminUsername, per.Read).Return(true, nil)
	suite.repository.On("Create", &result, nil)

	req, err := http.NewRequest(http.MethodPost, "/simulation", body)
	if err != nil {
		suite.Errorf(err, "Error creating HTTP Request.")
	}
	req.Header.Set("Content-Type", formWriter.FormDataContentType())

	suite.router.ServeHTTP(suite.recorder, req)

	var response Simulation
	b, err := ioutil.ReadAll(suite.recorder.Body)
	suite.NoError(err)
	suite.NoError(json.Unmarshal(b, &response))
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(uint(1), response.ID)
}
