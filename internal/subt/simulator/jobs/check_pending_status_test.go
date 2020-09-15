package jobs

import (
	"errors"
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

func TestCheckPendingStatus_ErrSimDoesNotExist(t *testing.T) {
	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation("", "", 0, nil, "")

	// Initialize fake simulation service
	svc := fake.NewService()
	err := errors.New("sim does not exist")
	svc.On("Get", gid).Return(sim, err)
	app := application.NewServices(svc)

	// Initialize job input and store
	input := state.NewStartSimulation(nil, app, gid)
	s := actions.NewStore(input)

	_, jobErr := CheckPendingStatus.Run(s, nil, nil, input)
	assert.Error(t, jobErr)
	assert.Equal(t, err, jobErr)
}
