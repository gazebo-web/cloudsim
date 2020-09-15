package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestCheckKind_Success(t *testing.T) {
	var state int
	s := actions.NewStore(&state)

	sim := fake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimParent, nil, "test")

	input := CheckKindInput{
		Simulation: sim,
		Kind:       simulations.SimParent,
	}

	result, err := CheckKind.Run(s, nil, &actions.Deployment{CurrentJob: "test"}, input)
	assert.NoError(t, err)

	output, ok := result.(CheckKindOutput)
	assert.True(t, ok)
	assert.Equal(t, input.Simulation, output)
}

func TestCheckKind_ErrWhenKindDoesNotMatch(t *testing.T) {
	var state int
	s := actions.NewStore(&state)

	sim := fake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimChild, nil, "test")

	input := CheckKindInput{
		Simulation: sim,
		Kind:       simulations.SimParent,
	}

	result, err := CheckKind.Run(s, nil, &actions.Deployment{CurrentJob: "test"}, input)
	assert.Error(t, err)
	assert.Equal(t, simulations.ErrIncorrectKind, err)
	assert.Nil(t, result)
}
