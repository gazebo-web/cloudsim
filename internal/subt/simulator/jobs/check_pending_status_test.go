package jobs

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	simctx "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
	"testing"
)

func TestCheckPendingStatus_Success(t *testing.T) {
	s := fake.NewService()

	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation(gid, simulations.StatusPending, simulations.SimSingle, nil, "test")

	s.On("Get", gid).Return(sim, nil)

	app := application.NewServices(s)

	ctx := context.Background()
	ctx = context.WithValue(ctx, simctx.CtxServices, app)
	ctx = actions.NewContext(ctx)

	input := &StartSimulationData{
		GroupID: gid,
	}

	result, err := CheckPendingStatus.Run(ctx, nil, nil, input)
	assert.NoError(t, err)

	output, ok := result.(*StartSimulationData)
	assert.True(t, ok)

	assert.Equal(t, input.GroupID, output.GroupID)

}

func TestCheckPendingStatus_ErrSimDoesNotExist(t *testing.T) {
	s := fake.NewService()

	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation("", "", 0, nil, "")

	err := errors.New("sim does not exist")
	s.On("Get", gid).Return(sim, err)

	app := application.NewServices(s)

	ctx := context.Background()
	ctx = context.WithValue(ctx, simctx.CtxServices, app)
	ctx = actions.NewContext(ctx)

	input := &StartSimulationData{
		GroupID: gid,
	}

	_, jobErr := CheckPendingStatus.Run(ctx, nil, nil, input)
	assert.Error(t, jobErr)
	assert.Equal(t, err, jobErr)
}
