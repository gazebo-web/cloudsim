package repository

import "errors"

var (
	// ErrNoFilter represents an error when no filter are provided.
	ErrNoFilter = errors.New("no filters provided")
	// ErrNoEntitiesUpdated represent an error when no entities were updated in the database
	// after an Update operation.
	ErrNoEntitiesUpdated = errors.New("no entities were updated")
	// ErrNoEntitiesDeleted represent an error when no entities were deleted in the database
	// after a Delete operation.
	ErrNoEntitiesDeleted = errors.New("no entities were deleted")
)

// Repository holds methods to CRUD an entity on a certain persistence layer.
type Repository interface {
	// Create is a bulk operation to create N amount of models passed in an array.
	Create([]Model) ([]Model, error)
	// Find writes on output the result of applying offset, limit and filters to find a range of models.
	Find(output interface{}, offset, limit *int, filters ...Filter) error
	// FindOne applies the given filters to write on entity the first result found.
	FindOne(entity Model, filters ...Filter) error
	// Update updates with the given data the different models that match filters.
	Update(data interface{}, filters ...Filter) error
	// Delete removes all the models that match filters.
	Delete(filters ...Filter) error
	// Model returns this repository's model.
	Model() Model
}
