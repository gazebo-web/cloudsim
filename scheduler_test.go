package main

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/scheduler"
	igntest "gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"context"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func createRegisterMockFn(wg *sync.WaitGroup, seconds int, circuit string) func(s *sim.Service, ctx context.Context, tx *gorm.DB) {
	return func(s *sim.Service, ctx context.Context, tx *gorm.DB) {
		scheduler.GetInstance().DoIn(
			func() {
				s.DeployHeldCircuitSimulations(ctx, tx, circuit)
				wg.Done()
			},
			seconds,
		)
	}
}

func createLaunchSimMockFn(launched *bool) func(s *sim.Service, ctx context.Context, tx *gorm.DB,
	dep *sim.SimulationDeployment) *ign.ErrMsg {
	return func(s *sim.Service, ctx context.Context, tx *gorm.DB, dep *sim.SimulationDeployment) *ign.ErrMsg {
		if em := s.GetApplications()[*dep.Application].ValidateSimulationLaunch(ctx, tx, dep); em != nil {
			return em
		}
		*launched = true
		return nil
	}
}

// Tests that the simulation scheduler adds all simulations for a circuit to the
// simulation queue once the competition date is met.
func TestRunScheduledJobOnCompetitionDay(t *testing.T) {
	// General test setup
	setup()

	// Target circuit
	circuit := "Urban Practice 2"

	// Setup a WaitGroup to pause the test until the scheduler thread fires
	wg := sync.WaitGroup{}
	wg.Add(1)

	restoreRegisterFn := ReplaceValue(
		&sim.RegisterSchedulableTasks,
		createRegisterMockFn(&wg, 2, circuit),
	)
	defer restoreRegisterFn()

	// Replace the LaunchSimulation function to avoid adding simulations to the queue
	launched := false
	restoreLaunchSimFn := ReplaceValue(
		&sim.LaunchSimulation,
		createLaunchSimMockFn(&launched),
	)
	defer restoreLaunchSimFn()

	// `RegisterSchedulableTasks` is called at application launch, which happens
	// before calling the `setup` function for this testcase. We are calling
	// `RegisterSchedulableTasks` again so that it runs AFTER calling `setup`,
	// allowing it to schedule launches using test circuit rules.
	// TODO: Having multiple tests call RegisterSchedulableTasks may cause
	//  unexpected issues.
	sim.RegisterSchedulableTasks(sim.SimServImpl.(*sim.Service), context.Background(), globals.Server.Db)

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
			testDesc:          "schedulerTest -- Launch simulation when reaching competition date",
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
		// Request a simulation before the competition date
		invokeURITestMultipartPOST(t, test.uriTest, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {})

		// Check that the simulation has not been launched
		assert.False(t, launched)

		// Wait for the scheduler to fire
		wg.Wait()

		// Check that the simulation has been launched
		assert.True(t, launched)
	})
}
