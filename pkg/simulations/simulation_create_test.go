package simulations

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/db"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
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

	db *gorm.DB

	userService   *users.ServiceMock
	admin         fuel.User
	user         fuel.User

	repository Repository
	service    Service
	controller Controller

	router   *mux.Router
	recorder *httptest.ResponseRecorder

	updateSimulation *SimulationUpdate
	createSimulation *SimulationCreate
	simulation       *Simulation

	uuid *uuid.MockUUID
}

func (suite *simulationTestSuite) SetupSuite() {

	suite.db = db.Must(db.NewDB(db.NewTestConfig()))
	suite.uuid = uuid.NewTestUUID()
	suite.userService = users.NewServiceMock()
	suite.repository = NewRepository(suite.db, tools.Sptr("platform_test"), tools.Sptr("app_test"))
	suite.service = NewService(NewServiceInput{
		Repository: suite.repository,
		Config:     ServiceConfig{
			Platform:    "platform_test",
			Application: "app_test",
			MaxDuration: time.Second,
		},
		UserService: suite.userService,
		UUID: suite.uuid,
	})

	suite.controller = NewController(NewControllerInput{
		Services: services{
			Simulation: suite.service,
			User:       suite.userService,
		},
		FormDecoder:       form.NewDecoder(),
		Validator:         validator.New(),
	})

	suite.router = mux.NewRouter()
	suite.recorder = httptest.NewRecorder()

	adminUsername := "root"
	suite.admin = fuel.User{
		Name:     tools.Sptr("Admin Root"),
		Username: &adminUsername,
		Email:    tools.Sptr("root@admin.com"),
	}

	userUsername := "test"
	suite.user = fuel.User{
		Name:     tools.Sptr("User Test"),
		Username: &userUsername,
		Email:    tools.Sptr("test@test.com"),
	}
}

func (suite *simulationTestSuite) BeforeTest(suiteName, testName string) {
	suite.db.AutoMigrate(&Simulation{})
}

func (suite *simulationTestSuite) AfterTest(suiteName, testName string) {
	suite.db.DropTableIfExists(&Simulation{})
}

func (suite *simulationTestSuite) TestCreate() {
	suite.router.HandleFunc("/simulation", func(writer http.ResponseWriter, request *http.Request) {
		result, err := suite.controller.Start(&suite.user, writer, request)
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
		"name":        "test1234",
		"circuit":		"test",
		"owner":       *suite.user.Username,
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

	var em *ign.ErrMsg
	returnedUUID := "aaaa-bbbb-cccc-dddd"
	suite.uuid.On("Generate").Return(returnedUUID)
	suite.userService.On("GetUserFromUsername", *suite.user.Username).Return(suite.user, em)
	suite.userService.On("IsSystemAdmin", *suite.user.Username).Return(false)
	// *sim.GroupID, []per.Action{per.Read, per.Write}, *sim.Owner, *createdSim.Application

	suite.userService.On("AddResourcePermission", *suite.user.Username, returnedUUID, per.Read).Return(true, em)
	suite.userService.On("AddResourcePermission", *suite.user.Username, returnedUUID, per.Write).Return(true, em)
	suite.userService.On("AddResourcePermission", "app_test", returnedUUID, per.Read).Return(true, em)
	suite.userService.On("AddResourcePermission", "app_test", returnedUUID, per.Write).Return(true, em)

	suite.userService.On("VerifyOwner", *suite.user.Username, *suite.user.Username, per.Read).Return(true, em)

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
}

func (suite *simulationTestSuite) TestCreateAdmin() {
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
		"name":        "test1234",
		"circuit":		"test",
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

	var em *ign.ErrMsg
	returnedUUID := "aaaa-bbbb-cccc-dddd"
	suite.uuid.On("Generate").Return(returnedUUID)
	suite.userService.On("GetUserFromUsername", *suite.admin.Username).Return(suite.admin, em)
	suite.userService.On("IsSystemAdmin", *suite.admin.Username).Return(true)

	suite.userService.On("AddResourcePermission", *suite.admin.Username, returnedUUID, per.Read).Return(true, em)
	suite.userService.On("AddResourcePermission", *suite.admin.Username, returnedUUID, per.Write).Return(true, em)
	suite.userService.On("AddResourcePermission", "app_test", returnedUUID, per.Read).Return(true, em)
	suite.userService.On("AddResourcePermission", "app_test", returnedUUID, per.Write).Return(true, em)

	suite.userService.On("VerifyOwner", *suite.admin.Username, *suite.admin.Username, per.Read).Return(true, em)

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
}