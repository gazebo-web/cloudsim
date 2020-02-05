package main

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	igntest "gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type launchHeldSimsTest struct {
	uriTest
}

func TestSystemAdminCanLaunchHeldSimulation(t *testing.T) {
	// General test setup
	setup()

	circuit := "Urban Circuit"
	createSimURI := "/1.0/simulations"
	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	createSubtForm := map[string]string{
		"name":        "sim1",
		"owner":       "TeamA",
		"circuit":     circuit,
		"robot_name":  "X1",
		"robot_type":  "X1_SENSOR_CONFIG_1",
		"robot_image": "infrastructureascode/aws-cli:latest",
	}

	var groupId string
	invokeURITestMultipartPOST(
		t,
		uriTest{
			testDesc:          "launchHeldSimulations -- Creating simulation deployment",
			URL:               createSimURI,
			jwtGen:            teamAUser1,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
		createSubtForm,
		func(bslice *[]byte, resp *igntest.AssertResponse) {
			dep := sim.SimulationDeployment{}
			json.Unmarshal(*bslice, &dep)
			groupId = *dep.GroupId
			assert.True(t, dep.Held)
		},
	)

	getURI := "/1.0/simulations/%s"
	launchHeldSimURI := "/1.0/simulations/%s/launch"
	sysAdmin := getDefaultTestJWT()
	test := launchHeldSimsTest{
		uriTest: uriTest{
			testDesc:          "launchHeldSimulations -- SystemAdmin can deploy a held simulation",
			URL:               fmt.Sprintf(launchHeldSimURI, groupId),
			jwtGen:            sysAdmin,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
	}

	t.Run(test.testDesc, func(t *testing.T) {
		invokeURITestMultipartPOST(t, test.uriTest, nil, func(bslice *[]byte, resp *igntest.AssertResponse) {
			dep := sim.SimulationDeployment{}
			require.NoError(t, json.Unmarshal(*bslice, &dep), "Unable to unmarshal response", string(*bslice))
			assert.False(t, dep.Held)
		})

		invokeURITest(t,
			uriTest{
				testDesc:          "launchHeldSimulations -- Get launched simulation deployment",
				URL:               fmt.Sprintf(getURI, groupId),
				jwtGen:            test.jwtGen,
				expErrMsg:         nil,
				ignoreErrorBody:   false,
				ignoreOptionsCall: false,
			},
			func(bslice *[]byte, resp *igntest.AssertResponse) {
				dep := sim.SimulationDeployment{}
				json.Unmarshal(*bslice, &dep)
				assert.False(t, dep.Held)
			},
		)
	})
}

func TestSystemAdminCanLaunchHeldMultisim(t *testing.T) {
	// General test setup
	setup()

	circuit := "Urban Practice 1"
	createSimURI := "/1.0/simulations"
	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	createSubtForm := map[string]string{
		"name":        "sim1",
		"owner":       "TeamA",
		"circuit":     circuit,
		"robot_name":  "X1",
		"robot_type":  "X1_SENSOR_CONFIG_1",
		"robot_image": "infrastructureascode/aws-cli:latest",
	}

	var groupId string
	invokeURITestMultipartPOST(
		t,
		uriTest{
			testDesc:          "launchHeldSimulations -- Creating simulation deployment",
			URL:               createSimURI,
			jwtGen:            teamAUser1,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
		createSubtForm,
		func(bslice *[]byte, resp *igntest.AssertResponse) {
			dep := sim.SimulationDeployment{}
			json.Unmarshal(*bslice, &dep)
			groupId = *dep.GroupId
			assert.True(t, dep.Held)
		},
	)

	getURI := "/1.0/simulations/%s"
	launchHeldSimURI := "/1.0/simulations/%s/launch"
	sysAdmin := getDefaultTestJWT()

	test := launchHeldSimsTest{
		uriTest: uriTest{
			testDesc:          "launchHeldSimulations -- SystemAdmin can deploy a held simulation",
			URL:               fmt.Sprintf(launchHeldSimURI, groupId),
			jwtGen:            sysAdmin,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
	}

	t.Run(test.testDesc, func(t *testing.T) {
		invokeURITestMultipartPOST(t, test.uriTest, nil, func(bslice *[]byte, resp *igntest.AssertResponse) {
			dep := sim.SimulationDeployment{}
			require.NoError(t, json.Unmarshal(*bslice, &dep), "Unable to unmarshal response", string(*bslice))
			assert.False(t, dep.Held)
		})

		invokeURITest(t,
			uriTest{
				testDesc:          "launchHeldSimulations -- Get launched simulation deployment",
				URL:               fmt.Sprintf(getURI, groupId),
				jwtGen:            test.jwtGen,
				expErrMsg:         nil,
				ignoreErrorBody:   false,
				ignoreOptionsCall: false,
			},
			func(bslice *[]byte, resp *igntest.AssertResponse) {
				dep := sim.SimulationDeployment{}
				json.Unmarshal(*bslice, &dep)
				assert.False(t, dep.Held)
			},
		)
	})
}

func TestUserCannotLaunchHeldSimulation(t *testing.T) {
	// General test setup
	setup()

	circuit := "Urban Circuit"
	createSimURI := "/1.0/simulations"
	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	createSubtForm := map[string]string{
		"name":        "sim1",
		"owner":       "TeamA",
		"circuit":     circuit,
		"robot_name":  "X1",
		"robot_type":  "X1_SENSOR_CONFIG_1",
		"robot_image": "infrastructureascode/aws-cli:latest",
	}

	var dep sim.SimulationDeployment
	invokeURITestMultipartPOST(
		t,
		uriTest{
			testDesc:          "launchHeldSimulations -- Creating simulation deployment",
			URL:               createSimURI,
			jwtGen:            teamAUser1,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
		createSubtForm,
		func(bslice *[]byte, resp *igntest.AssertResponse) {
			json.Unmarshal(*bslice, &dep)
		},
	)

	launchHeldSimURI := "/1.0/simulations/%s/launch"
	test := launchHeldSimsTest{
		uriTest: uriTest{
			testDesc:          "launchHeldSimulations -- User cannot deploy a held simulation",
			URL:               fmt.Sprintf(launchHeldSimURI, *dep.GroupId),
			jwtGen:            teamAUser1,
			expErrMsg:         ign.NewErrorMessage(ign.ErrorUnauthorized),
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
	}
	t.Run(test.testDesc, func(t *testing.T) {
		invokeURITestMultipartPOST(t, test.uriTest, nil, func(bslice *[]byte, resp *igntest.AssertResponse) {})
	})
}
