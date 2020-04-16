package deprecated

import (
	igntest "gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTeamIsQualified(t *testing.T) {
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

	test := createSimulationTest{
		uriTest: uriTest{
			testDesc:          "teamQualification -- Team is qualified to participate in Urban Circuit",
			URL:               createSimURI,
			jwtGen:            teamAUser1,
			expErrMsg:         nil,
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
		circuit:    circuit,
		owner:      "TeamA",
		robotName:  nil,
		expCreator: "TeamAUser1",
	}

	t.Run(test.testDesc, func(t *testing.T) {
		invokeURITestMultipartPOST(t, test.uriTest, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {
			dep := sim.SimulationDeployment{}
			require.NoError(t, json.Unmarshal(*bslice, &dep), "Unable to unmarshal response", string(*bslice))
		})
	})
}

func TestTeamIsNotQualified(t *testing.T) {
	// General test setup
	setup()

	circuit := "Urban Circuit"
	createSimURI := "/1.0/simulations"
	teamBUser1 := newJWT(createJWTForIdentity(t, "TeamBUser1"))

	createSubtForm := map[string]string{
		"name":        "sim2",
		"owner":       "TeamB",
		"circuit":     circuit,
		"robot_name":  "X1",
		"robot_type":  "X1_SENSOR_CONFIG_1",
		"robot_image": "infrastructureascode/aws-cli:latest",
	}

	test := createSimulationTest{
		uriTest: uriTest{
			testDesc:          "teamQualification -- Team is not qualified to participate in Urban Circuit",
			URL:               createSimURI,
			jwtGen:            teamBUser1,
			expErrMsg:         sim.NewErrorMessage(sim.ErrorNotQualified),
			ignoreErrorBody:   false,
			ignoreOptionsCall: false,
		},
		circuit:    circuit,
		owner:      "TeamB",
		robotName:  nil,
		expCreator: "TeamBUser1",
	}

	t.Run(test.testDesc, func(t *testing.T) {
		invokeURITestMultipartPOST(t, test.uriTest, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {
			dep := sim.SimulationDeployment{}
			require.NoError(t, json.Unmarshal(*bslice, &dep), "Unable to unmarshal response", string(*bslice))
		})
	})
}
