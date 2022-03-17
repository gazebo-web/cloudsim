package repository

// Repository holds methods to CRUD an entity on a certain persistence layer.
type Repository interface {
	// Create saves the given input entity in a persistence layer.
	// The result is written back on the input.
	Create(in interface{}) error
	// Get writes on out the data of the entity identified by id.
	Get(out interface{}, id uint) error
	// Update updates an entity identified by id with the input data.
	// The result is written back on the input.
	Update(in interface{}, id uint) error
	// Remove removes the entity identified by id and writes on out the data of the element that has been removed.
	Remove(out interface{}, id uint) error
}
