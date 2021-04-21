package jobs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
	"time"
)

func TestCheckSimulationNoErrors_Success(t *testing.T) {

	input := CheckSimulationNoErrorInput([]simulations.Simulation{
		fake.NewSimulation(
			"aaaa-bbbb-cccc-dddd", simulations.StatusPending, simulations.SimSingle, nil, "", time.Minute,
		),
		fake.NewSimulation(
			"eeee-ffff-gggg-hhhh", simulations.StatusPending, simulations.SimSingle, nil, "", time.Minute,
		),
		fake.NewSimulation(
			"iiii-jjjj-kkkk-llll", simulations.StatusPending, simulations.SimSingle, nil, "", time.Minute,
		),
	})

	result, err := CheckSimulationNoError.Run(nil, nil, nil, input)

	output, ok := result.(CheckSimulationNoErrorOutput)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.NoError(t, output.Error)
}

func TestCheckSimulationNoErrors_WithError(t *testing.T) {

	simulationError := simulations.Error("test")
	input := CheckSimulationNoErrorInput([]simulations.Simulation{
		fake.NewSimulation(
			"aaaa-bbbb-cccc-dddd", simulations.StatusPending, simulations.SimSingle, nil, "", time.Minute,
		),
		fake.NewSimulation(
			"eeee-ffff-gggg-hhhh", simulations.StatusPending, simulations.SimSingle, nil, "", time.Minute,
		),
		fake.NewSimulation(
			"iiii-jjjj-kkkk-llll", simulations.StatusPending, simulations.SimSingle, &simulationError, "", time.Minute,
		),
	})

	result, err := CheckSimulationNoError.Run(nil, nil, nil, input)

	output, ok := result.(CheckSimulationNoErrorOutput)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Error(t, output.Error)

	// The actual error should have information from the simulation that returned the error.
	assert.Equal(t,
		fmt.Sprintf("simulation [%s] with error status [%s]", "iiii-jjjj-kkkk-llll", simulationError),
		output.Error.Error(),
	)
}
