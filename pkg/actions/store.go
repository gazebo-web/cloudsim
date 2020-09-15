package actions

// State represents the store states. It's used to save data, artifacts and dependencies that will be used across jobs.
type State interface{}

// Store groups a set of methods to interact with the state.
type Store interface {
	// State returns the store's state.
	State() State

	// SetState changes the current state with the given state.
	SetState(state State)
}

// store is a Store implementation.
type store struct {
	// state holds the store's state.
	state State
}

func (s *store) SetState(state State) {
	s.state = state
}

// State returns the store's state.
func (s *store) State() State {
	return s.state
}

// NewStore initializes a new store with the given initial state.
func NewStore(initialState State) Store {
	if initialState == nil {
		panic("no initial state")
	}
	return &store{
		state: initialState,
	}
}
