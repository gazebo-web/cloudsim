package actions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateStore_PanicsWhenInitialStateIsNil(t *testing.T) {
	assert.Panics(t, func() {
		s := CreateStore(nil, nil)
		assert.Nil(t, s)
	})
}

func TestCreateStore_SuccessWhenInitialStateIsProvided(t *testing.T) {
	type someState struct {
		value int
	}

	assert.NotPanics(t, func() {
		s := CreateStore(&someState{value: 1}, nil)
		assert.NotNil(t, s)
	})
}

func TestDispatchMutationFailsWhenKeyIsZeroValue(t *testing.T) {
	type someState struct {
		value int
	}
	s := CreateStore(&someState{value: 1}, nil)
	err := s.Dispatch("", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrStoreDispatchNilKey, err)
}

func TestDispatchMutationFailsWhenMutationDoesNotExist(t *testing.T) {
	type someState struct {
		value int
	}
	s := CreateStore(&someState{value: 1}, nil)
	err := s.Dispatch("test", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrStoreInvalidMutation, err)
}

func TestDispatchMutationSuccessWhenMutationExists(t *testing.T) {
	type someState struct {
		value int
	}

	m := Mutation(func(state State, value interface{}) error {
		input := state.(*someState)
		increment := value.(int)
		input.value += increment
		return nil
	})

	var storeState someState

	s := CreateStore(&storeState, map[string]Mutation{
		"adder": m,
	})

	assert.Equal(t, 0, storeState.value)
	err := s.Dispatch("adder", 5)
	assert.NoError(t, err)
	assert.Equal(t, 5, storeState.value)
}
