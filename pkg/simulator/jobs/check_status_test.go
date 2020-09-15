package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestCheckStatus_Success(t *testing.T) {
	var state int
	s := actions.NewStore(&state)

	sim := fake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimChild, nil, "test")

	input := CheckStatusInput{
		Simulation: sim,
		Status:     simulations.StatusPending,
	}

	result, err := CheckStatus.Run(s, nil, &actions.Deployment{CurrentJob: "test"}, input)
	assert.NoError(t, err)

	output, ok := result.(CheckStatusOutput)
	assert.True(t, ok)
	assert.Equal(t, input.Simulation, output)
}

func TestCheckStatus_ErrWhenStatusDoesNotMatch(t *testing.T) {
	var state int
	s := actions.NewStore(&state)

	sim := fake.NewSimulation("test-group-id", simulations.StatusRunning, simulations.SimChild, nil, "test")

	input := CheckStatusInput{
		Simulation: sim,
		Status:     simulations.StatusPending,
	}

	result, err := CheckStatus.Run(s, nil, &actions.Deployment{CurrentJob: "test"}, input)
	assert.Error(t, err)
	assert.Equal(t, simulations.ErrIncorrectStatus, err)
	assert.Nil(t, result)
}
