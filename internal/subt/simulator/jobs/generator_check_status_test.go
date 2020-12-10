package jobs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestCheckSimulationStatus_Success(t *testing.T) {
	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation(gid, simulations.StatusPending, simulations.SimSingle, nil, "test")

	// Initialize fake simulation service
	svc := fake.NewService()
	svc.On("Get", gid).Return(sim, nil)
	app := application.NewServices(svc)

	tracksService := tracks.NewService(nil, nil, nil)

	subt := subtapp.NewServices(app, tracksService)

	// Initialize job input and store
	input := state.NewStartSimulation(nil, subt, gid)
	s := actions.NewStore(input)

	job := GenerateCheckStatusJob(CheckStatusJobConfig{
		Name:       "test",
		Status:     simulations.StatusPending,
		InputType:  &state.StartSimulation{},
		OutputType: &state.StartSimulation{},
	})

	result, err := job.Run(s, nil, nil, input)
	require.NoError(t, err)

	output, ok := result.(*state.StartSimulation)
	assert.True(t, ok)

	assert.Equal(t, input.GroupID, output.GroupID)

}

func TestCheckSimulationStatus_ErrSimInvaludStatus(t *testing.T) {
	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation(gid, simulations.StatusRunning, simulations.SimSingle, nil, "test")

	// Initialize fake simulation service
	svc := fake.NewService()
	svc.On("Get", gid).Return(sim, nil)
	app := application.NewServices(svc)

	tracksService := tracks.NewService(nil, nil, nil)

	subt := subtapp.NewServices(app, tracksService)

	// Initialize job input and store
	input := state.NewStartSimulation(nil, subt, gid)
	s := actions.NewStore(input)

	job := GenerateCheckStatusJob(CheckStatusJobConfig{
		Name:       "test",
		Status:     simulations.StatusPending,
		InputType:  &state.StartSimulation{},
		OutputType: &state.StartSimulation{},
	})

	_, err := job.Run(s, nil, nil, input)
	assert.Error(t, err)
	assert.Equal(t, simulations.ErrIncorrectStatus, err)
}
