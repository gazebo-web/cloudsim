package jobs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestGenerateSetSimulationStatusJob(t *testing.T) {
	// Initialize job to change status to running
	job := GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
		Name:       "test",
		Status:     simulations.StatusRunning,
		InputType:  &state.StartSimulation{},
		OutputType: &state.StartSimulation{},
		PreHooks:   []actions.JobFunc{setStartState, returnGroupIDFromStartState},
		PostHooks:  nil,
	})

	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation(gid, simulations.StatusPending, simulations.SimSingle, nil, "test")

	// Initialize fake simulation service
	svc := fake.NewService()
	svc.On("UpdateStatus", gid, simulations.StatusRunning).Return(error(nil)).Run(func(args mock.Arguments) {
		sim.SetStatus(simulations.StatusRunning)
	})
	app := application.NewServices(svc, nil)

	tracksService := tracks.NewService(nil, nil, nil)

	subt := subtapp.NewServices(app, tracksService, nil)

	// Initialize store
	initialState := state.NewStartSimulation(nil, subt, gid)
	s := actions.NewStore(&initialState)

	_, err := job.Run(s, nil, nil, initialState)

	assert.NoError(t, err)
	assert.Equal(t, simulations.StatusRunning, sim.GetStatus())

}
