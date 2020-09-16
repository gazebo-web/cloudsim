package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestCheckPendingStatus_Success(t *testing.T) {
	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation(gid, simulations.StatusPending, simulations.SimSingle, nil, "test")

	// Initialize fake simulation service
	svc := fake.NewService()
	svc.On("Get", gid).Return(sim, nil)
	app := application.NewServices(svc)

	// Initialize job input and store
	input := state.NewStartSimulation(nil, app, gid)
	s := actions.NewStore(input)

	result, err := CheckPendingStatus.Run(s, nil, nil, input)
	assert.NoError(t, err)

	output, ok := result.(*state.StartSimulation)
	assert.True(t, ok)

	assert.Equal(t, input.GroupID, output.GroupID)

}

func TestCheckPendingStatus_ErrSimNotPending(t *testing.T) {
	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation(gid, simulations.StatusRunning, simulations.SimSingle, nil, "test")

	// Initialize fake simulation service
	svc := fake.NewService()
	svc.On("Get", gid).Return(sim, nil)
	app := application.NewServices(svc)

	// Initialize job input and store
	input := state.NewStartSimulation(nil, app, gid)
	s := actions.NewStore(input)

	_, err := CheckPendingStatus.Run(s, nil, nil, input)
	assert.Error(t, err)
	assert.Equal(t, simulations.ErrIncorrectStatus, err)
}
