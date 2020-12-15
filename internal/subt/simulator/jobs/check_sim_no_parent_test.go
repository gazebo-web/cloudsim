package jobs

import (
	"github.com/stretchr/testify/assert"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	tfake "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestCheckSimIsParent(t *testing.T) {
	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation(gid, simulations.StatusPending, simulations.SimSingle, nil, "test")

	// Initialize fake simulation service
	svc := fake.NewService()
	svc.On("Get", gid).Return(sim, nil)

	// Initialize tracks service
	trackService := tfake.NewService()

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(svc), trackService)

	// Initialize job input and store
	input := state.NewStartSimulation(nil, app, gid)
	s := actions.NewStore(input)

	result, err := CheckSimulationIsNotParent.Run(s, nil, nil, input)
	assert.NoError(t, err)

	output, ok := result.(*state.StartSimulation)
	assert.True(t, ok)

	assert.Equal(t, input.GroupID, output.GroupID)

}

func TestCheckSimIsParent_ErrSimIsParent(t *testing.T) {
	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := fake.NewSimulation(gid, simulations.StatusPending, simulations.SimParent, nil, "test")

	// Initialize fake simulation service
	svc := fake.NewService()
	svc.On("Get", gid).Return(sim, nil)

	// Initialize tracks service
	trackService := tfake.NewService()

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(svc), trackService)

	// Initialize job input and store
	input := state.NewStartSimulation(nil, app, gid)
	s := actions.NewStore(input)

	_, err := CheckSimulationIsNotParent.Run(s, nil, nil, input)
	assert.Error(t, err)
	assert.Equal(t, simulations.ErrIncorrectKind, err)
}
