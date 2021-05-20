package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"testing"
)

func TestReturnState(t *testing.T) {
	initValue := 15
	s := actions.NewStore(&initValue)
	value, err := returnState(s, nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, &initValue, value)
}

func TestSetStartState(t *testing.T) {
	// Initialize empty store
	s := actions.NewStore(&state.StartSimulation{})

	// Set start state
	startState := state.StartSimulation{GroupID: "test"}

	value, err := setStartState(s, nil, nil, &startState)
	assert.NoError(t, err)

	// Check that the returned value is the same as the one we passed.
	parsed, ok := value.(*state.StartSimulation)
	assert.True(t, ok)
	assert.Equal(t, &startState, parsed)
}

func TestReturnGroupIDFromStartState(t *testing.T) {
	s := actions.NewStore(&state.StartSimulation{
		GroupID: "test",
	})

	value, err := returnGroupIDFromStartState(s, nil, nil, nil)
	assert.NoError(t, err)

	parsed, ok := value.(simulations.GroupID)
	assert.True(t, ok)
	assert.Equal(t, simulations.GroupID("test"), parsed)
}
