package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
	"time"
)

func TestCheckSimulationStatus_Success(t *testing.T) {
	var state int
	s := actions.NewStore(&state)

	sim := fake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimChild, nil, "test", 1*time.Minute, nil, nil)

	input := CheckSimulationStatusInput{
		Simulation: sim,
		Status:     simulations.StatusPending,
	}

	result, err := CheckSimulationStatus.Run(s, nil, &actions.Deployment{CurrentJob: "test"}, input)
	assert.NoError(t, err)

	output, ok := result.(CheckSimulationStatusOutput)
	assert.True(t, ok)
	assert.True(t, bool(output))
}

func TestCheckSimulationStatus_ErrWhenStatusDoesNotMatch(t *testing.T) {
	var state int
	s := actions.NewStore(&state)

	sim := fake.NewSimulation("test-group-id", simulations.StatusRunning, simulations.SimChild, nil, "test", 1*time.Minute, nil, nil)

	input := CheckSimulationStatusInput{
		Simulation: sim,
		Status:     simulations.StatusPending,
	}

	result, err := CheckSimulationStatus.Run(s, nil, &actions.Deployment{CurrentJob: "test"}, input)
	assert.NoError(t, err)

	output, ok := result.(CheckSimulationStatusOutput)
	assert.True(t, ok)
	assert.False(t, bool(output))
}
