package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	igntest "gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	"testing"
)

func TestWebsocketAddressUser(t *testing.T) {
	// General test setup
	setup()

	circuit := "Virtual Stix"
	createSimURI := "/1.0/simulations"
	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))
	teamBUser1 := newJWT(createJWTForIdentity(t, "TeamBUser1"))

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
		t,
		uriTest{
			testDesc:          "WebSocket Address Test -- Creating simulation deployment",
			URL:               createSimURI,
			jwtGen:            teamAUser1,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
		createSubtForm,
		func(bslice *[]byte, resp *igntest.AssertResponse) {
			var dep sim.SimulationDeployment
			assert.NoError(t, json.Unmarshal(*bslice, &dep))
			groupID = *dep.GroupID
			assert.True(t, dep.Held)
		},
	)

	websocketAddr := "/1.0/simulations/%s/websocket"
	testA := uriTest{
		testDesc:        "WebSocket Address Test -- User TeamA should get websocket address",
		URL:             fmt.Sprintf(websocketAddr, groupID),
		jwtGen:          teamAUser1,
		expErrMsg:       nil,
		ignoreErrorBody: false,
	}

	t.Run(testA.testDesc, func(t *testing.T) {
		invokeURITest(t, testA, func(bslice *[]byte, resp *igntest.AssertResponse) {
			var wsResp sim.WebsocketAddressResponse
			assert.NoError(t, json.Unmarshal(*bslice, &wsResp))
			assert.True(t, sim.IsWebsocketAddress(wsResp.Address, &groupID))
		})
	})

	testB := uriTest{
		testDesc:        "WebSocket Address Test -- User TeamB shouldn't get websocket address from TeamA",
		URL:             fmt.Sprintf(websocketAddr, groupID),
		jwtGen:          teamBUser1,
		expErrMsg:       ign.NewErrorMessage(ign.ErrorUnauthorized),
		ignoreErrorBody: false,
	}

	t.Run(testB.testDesc, func(t *testing.T) {
		invokeURITest(t, testB, func(bslice *[]byte, resp *igntest.AssertResponse) {
			var wsResp sim.WebsocketAddressResponse
			assert.NoError(t, json.Unmarshal(*bslice, &wsResp))
			assert.NotEmpty(t, wsResp.Token)
		})
	})
}

func TestWebsocketAddressAdmin(t *testing.T) {
	// General test setup
	setup()

	circuit := "Virtual Stix"
	createSimURI := "/1.0/simulations"
	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))
	sysAdmin := getDefaultTestJWT()

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
		t,
		uriTest{
			testDesc:          "WebSocket Address Test -- Creating simulation deployment",
			URL:               createSimURI,
			jwtGen:            teamAUser1,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
		createSubtForm,
		func(bslice *[]byte, resp *igntest.AssertResponse) {
			var dep sim.SimulationDeployment
			assert.NoError(t, json.Unmarshal(*bslice, &dep))
			groupID = *dep.GroupID
			assert.True(t, dep.Held)
		},
	)

	websocketAddr := "/1.0/simulations/%s/websocket"
	testA := uriTest{
		testDesc:        "WebSocket Address Test -- Admin should get websocket address for every simulation",
		URL:             fmt.Sprintf(websocketAddr, groupID),
		jwtGen:          sysAdmin,
		expErrMsg:       nil,
		ignoreErrorBody: false,
	}

	t.Run(testA.testDesc, func(t *testing.T) {
		invokeURITest(t, testA, func(bslice *[]byte, resp *igntest.AssertResponse) {
			var wsResp sim.WebsocketAddressResponse
			assert.NoError(t, json.Unmarshal(*bslice, &wsResp))
			assert.True(t, sim.IsWebsocketAddress(wsResp.Address, &groupID))
			assert.NotEmpty(t, wsResp.Token)
		})
	})
}
