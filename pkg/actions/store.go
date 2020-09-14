package actions

import (
	"errors"
)

var (
	// ErrStoreDispatchNilKey is triggered when the Dispatch method is called with a nil key.
	ErrStoreDispatchNilKey = errors.New("dispatch called with nil key")
	// ErrStoreInvalidMutation is triggered when the Dispatch method calls an nonexistent mutation.
	ErrStoreInvalidMutation = errors.New("mutation does not exist")
)

// State represents the store states. It's used to save data, artifacts and dependencies that will be used across jobs.
type State interface{}

// Mutation represents a change in the store's state. It returns an error if performing the change fails.
type Mutation func(state State, value interface{}) error

// Store groups a set of methods to interact with the state and trigger mutations on it.
type Store interface {
	// State returns the store's state.
	State() State
	// Dispatch dispatches a mutation identified by the given key.
	// It passes the given value as the mutation's input.
	Dispatch(key string, value interface{}) error
}

// store is a Store implementation.
type store struct {
	// state holds the store's state.
	state State
	// mutations holds the reference to every mutation available on this store.
	mutations map[string]Mutation
}

// State returns the store's state.
func (s *store) State() State {
	panic("implement me")
}

// Dispatch dispatches a mutation identified by the given key.
func (s *store) Dispatch(key string, value interface{}) error {
	if key == "" {
		return ErrStoreDispatchNilKey
	}
	fn, ok := s.mutations[key]
	if !ok {
		return ErrStoreInvalidMutation
	}
	return fn(s.state, value)
}

// CreateStore initializes a new store with the given initial state and the given mutations.
// The initialState argument must be provided in order to guarantee the store's underlying state data type.
// If no initialState is provided, it will panic.
func CreateStore(initialState State, mutations map[string]Mutation) Store {
	if initialState == nil {
		panic("no initial state")
	}
	return &store{
		state:     initialState,
		mutations: mutations,
	}
}
