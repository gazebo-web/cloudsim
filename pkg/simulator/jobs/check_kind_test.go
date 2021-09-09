package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
	"time"
)

func TestCheckSimulationKind_Success(t *testing.T) {
	var state int
	s := actions.NewStore(&state)

	sim := fake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimParent, nil, "test", 1*time.Minute, nil)

	input := CheckSimulationKindInput{
		Simulation: sim,
		Kind:       simulations.SimParent,
	}

	result, err := CheckSimulationKind.Run(s, nil, &actions.Deployment{CurrentJob: "test"}, input)
	assert.NoError(t, err)

	output, ok := result.(CheckSimulationKindOutput)
	assert.True(t, ok)
	assert.True(t, bool(output))
}

func TestCheckSimulationKind_ReturnsFalseWhenKindDoesNotMatch(t *testing.T) {
	var state int
	s := actions.NewStore(&state)

	sim := fake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimChild, nil, "test", 1*time.Minute, nil)

	input := CheckSimulationKindInput{
		Simulation: sim,
		Kind:       simulations.SimParent,
	}

	result, err := CheckSimulationKind.Run(s, nil, &actions.Deployment{CurrentJob: "test"}, input)
	assert.NoError(t, err)
	output, ok := result.(CheckSimulationKindOutput)
	assert.True(t, ok)
	assert.False(t, bool(output))
}
