package actions

// Store is in charge of loading and saving information across jobs for a single action.
// It includes the respective methods to get and set the information to and from a running job.
type Store interface {
	// Get gets the data from the action.
	Get() interface{}
	// Set sets the data in the action.
	Set(interface{}) error
	// Load loads the data from a persistent storage into the store.
	Load() error
	// Save saves the data from the store into a persistent storage.
	Save() error
}

// store is a Store implementation.
type store struct{}

// Get gets the data from the action.
func (s store) Get() interface{} {
	panic("implement me")
}

// Set sets the data in the action.
func (s store) Set(i interface{}) error {
	panic("implement me")
}

// Load loads the data from a persistent storage into the store.
func (s store) Load() error {
	panic("implement me")
}

// Save saves the data from the store into a persistent storage.
func (s store) Save() error {
	panic("implement me")
}

// NewStore initializes a new Store using the default store.
func NewStore() Store {
	return &store{}
}
