package jobs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestSetStoppedAt(t *testing.T) {
	simservice := fake.NewService()
	base := application.NewServices(simservice, nil, nil)
	app := subtapp.NewServices(base, nil, nil)

	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	simservice.On("MarkStopped", gid).Return(error(nil))

	initialState := state.NewStopSimulation(nil, app, gid)
	s := actions.NewStore(initialState)

	_, err := SetStoppedAt.Run(s, nil, nil, initialState)

	assert.NoError(t, err)
}

func TestSetStoppedAtWithError(t *testing.T) {
	simservice := fake.NewService()
	base := application.NewServices(simservice, nil, nil)
	app := subtapp.NewServices(base, nil, nil)

	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	simservice.On("MarkStopped", gid).Return(errors.New("test error"))

	initialState := state.NewStopSimulation(nil, app, gid)
	s := actions.NewStore(initialState)

	_, err := SetStoppedAt.Run(s, nil, nil, initialState)

	assert.Error(t, err)
}
