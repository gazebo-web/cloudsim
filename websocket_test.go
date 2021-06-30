package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	igntest "gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	"sync"
	"testing"
	"time"
)

func TestWebsocketAddressSuite(t *testing.T) {
	suite.Run(t, &WebsocketAddressTestSuite{})
}

type WebsocketAddressTestSuite struct {
	suite.Suite
	singleSimCircuit string
	singleSimGroupID string
	multiSimCircuit  string
	multiSimGroupID  string
	createSimURI     string
	websocketURI     string
	jwtTeamAUser1    *testJWT
	jwtTeamBUser1    *testJWT
	jwtSysAdmin      *testJWT
}

func (suite *WebsocketAddressTestSuite) SetupSuite() {
	suite.singleSimCircuit = "Virtual Stix"
	suite.multiSimCircuit = "Urban Practice 3"
	suite.createSimURI = "/1.0/simulations"
	suite.websocketURI = "/1.0/simulations/%s/websocket"
	suite.jwtTeamAUser1 = newJWT(createJWTForIdentity(suite.T(), "TeamAUser1"))
	suite.jwtTeamBUser1 = newJWT(createJWTForIdentity(suite.T(), "TeamBUser1"))
	suite.jwtSysAdmin = getDefaultTestJWT()

	// Initialize the database. Note that this is only performed once for the entire suite.
	setup()

	// Create simulations.
	suite.singleSimGroupID = suite.setSimulationToRunning(suite.requestSimulation(suite.singleSimCircuit, suite.jwtTeamAUser1))
	suite.multiSimGroupID = suite.setSimulationToRunning(suite.getChildSimGroupID(
		suite.setSimulationToRunning(suite.requestSimulation(suite.multiSimCircuit, suite.jwtTeamAUser1)),
	))
}

// requestSimulation requests a simulation launch
func (suite *WebsocketAddressTestSuite) requestSimulation(circuit string, ownerJWT *testJWT) string {

	// A WaitGroup is set to wait until the async worker finishes launching the simulation
	var wg sync.WaitGroup
	if circuit == suite.singleSimCircuit {
		wg.Add(1)
	} else if circuit == suite.multiSimCircuit {
		// HACK This is very tightly coupled to the number of children set in the test circuit
		wg.Add(6)
	} else {
		suite.Fail("Please use one of the circuits in this test suite.")
	}

	// Use the notify signal to mark the worker as done
	cb := func(poolEvent sim.PoolEvent, groupID string, result interface{}, err error) {
		if poolEvent == sim.PoolStartSimulation {
			wg.Done()
		}
	}
	sim.SimServImpl.(*sim.Service).SetPoolEventsListener(cb)

	createSubtForm := map[string]string{
		"name":        "sim1",
		"owner":       "TeamA",
		"circuit":     circuit,
		"robot_name":  "X1",
		"robot_type":  "X1_SENSOR_CONFIG_1",
		"robot_image": "infrastructureascode/aws-cli:latest",
	}

	var groupID string

	invokeURITestMultipartPOST(
		suite.T(),
		uriTest{
			testDesc:          "WebSocket Address Test -- Creating simulation deployment",
			URL:               suite.createSimURI,
			jwtGen:            ownerJWT,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
		createSubtForm,
		func(bslice *[]byte, resp *igntest.AssertResponse) {
			var dep sim.SimulationDeployment
			assert.NoError(suite.T(), json.Unmarshal(*bslice, &dep))
			groupID = *dep.GroupID
		},
	)

	// Wait up to 5 seconds for the async worker to finish
	suite.False(wgTimeoutWait(&wg, 5*time.Second))

	return groupID
}

// testWebsocketAddress tests that a specific user has access to a simulation's websocket address.
func (suite *WebsocketAddressTestSuite) testWebsocketAddress(testDesc string, groupID string, requestJWT *testJWT,
	expErrMsg *ign.ErrMsg) {

	testURI := uriTest{
		testDesc:        testDesc,
		URL:             fmt.Sprintf(suite.websocketURI, groupID),
		jwtGen:          requestJWT,
		expErrMsg:       expErrMsg,
		ignoreErrorBody: false,
	}

	suite.Run(testURI.testDesc, func() {
		invokeURITest(suite.T(), testURI, func(bslice *[]byte, resp *igntest.AssertResponse) {
			var wsResp sim.WebsocketAddressResponse
			suite.NoError(json.Unmarshal(*bslice, &wsResp))
			// Validate the address structure
			suite.True(sim.IsWebsocketAddress(wsResp.Address, &groupID))
			// Validate that a token was included
			suite.NotEmpty(wsResp.Token)
		})
	})
}

func (suite *WebsocketAddressTestSuite) getChildSimGroupID(groupID string) string {
	return fmt.Sprintf("%s-c-1", groupID)
}

func (suite *WebsocketAddressTestSuite) setSimulationToRunning(groupID string) string {
	// Update status to running
	suite.Require().NoError(sim.SimServImpl.(*sim.Service).ServiceAdaptor.UpdateStatus(
		simulations.GroupID(groupID),
		simulations.StatusRunning),
	)

	return groupID
}

func (suite *WebsocketAddressTestSuite) TestWebsocketAddressUser() {

	suite.testWebsocketAddress(
		"WebSocket Address Test -- User TeamA should get websocket address",
		suite.singleSimGroupID,
		suite.jwtTeamAUser1,
		nil,
	)

	suite.testWebsocketAddress(
		"WebSocket Address Test -- User TeamB shouldn't get websocket address from TeamA",
		suite.singleSimGroupID,
		suite.jwtTeamBUser1,
		ign.NewErrorMessage(ign.ErrorUnauthorized),
	)
}

func (suite *WebsocketAddressTestSuite) TestWebsocketAddressAdmin() {

	suite.testWebsocketAddress(
		"WebSocket Address Test -- Admin should get websocket address for every simulation",
		suite.singleSimGroupID,
		suite.jwtSysAdmin,
		nil,
	)
}

func (suite *WebsocketAddressTestSuite) TestWebsocketAddressChildSimulations() {

	suite.testWebsocketAddress(
		"WebSocket Address Test -- User TeamA should not get websocket address for child simulations",
		suite.multiSimGroupID,
		suite.jwtTeamAUser1,
		ign.NewErrorMessage(ign.ErrorUnauthorized),
	)

	suite.testWebsocketAddress(
		"WebSocket Address Test -- Admin should get websocket address for child simulations",
		suite.multiSimGroupID,
		suite.jwtSysAdmin,
		nil,
	)
}
