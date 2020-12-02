package jobs

import (
	"github.com/stretchr/testify/assert"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestCheckSimulationPendingStatus_Success(t *testing.T) {
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

	result, err := CheckSimulationPendingStatus.Run(s, nil, nil, input)
	assert.NoError(t, err)

	output, ok := result.(*state.StartSimulation)
	assert.True(t, ok)

	assert.Equal(t, input.GroupID, output.GroupID)

}

func TestCheckSimulationPendingStatus_ErrSimNotPending(t *testing.T) {
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

	_, err := CheckSimulationPendingStatus.Run(s, nil, nil, input)
	assert.Error(t, err)
	assert.Equal(t, simulations.ErrIncorrectStatus, err)
}
