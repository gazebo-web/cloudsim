package main

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	igntest "gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"context"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockLaunchSimulation(launched *bool) func(s *sim.Service, ctx context.Context, tx *gorm.DB, dep *sim.SimulationDeployment) *ign.ErrMsg {
	return func(s *sim.Service, ctx context.Context, tx *gorm.DB, dep *sim.SimulationDeployment) *ign.ErrMsg {
		if em := s.GetApplications()[*dep.Application].ValidateSimulationLaunch(ctx, tx, dep); em != nil {
			*launched = false
			return em
		}
		return nil
	}
}

func TestLaunchSimulationBeforeCompetitionDay(t *testing.T) {
	// General test setup
	setup()

	launched := true

	restoreFn := ReplaceValue(&sim.LaunchSimulation, mockLaunchSimulation(&launched))
	defer restoreFn()

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
			testDesc:          "holdSubmission -- Don't launch simulation before competition day",
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
			assert.Equal(t, false, launched)
		})

	})
}

func TestLaunchSimulationOnCompetitionDay(t *testing.T) {
	// General test setup
	setup()

	launched := true

	restoreFn := ReplaceValue(&sim.LaunchSimulation, mockLaunchSimulation(&launched))
	defer restoreFn()

	circuit := "Urban Practice 3"
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
			testDesc:          "holdSubmission -- Launch simulation after competition day",
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
			assert.Equal(t, true, launched)
		})
	})
}
