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
	Create([]Model) ([]Model, error)
	Find(output interface{}, offset, limit *int, filters ...Filter) error
	FindOne(entity Model, filters ...Filter) error
	Update(data interface{}, filters ...Filter) error
	Delete(filters ...Filter) error
	Model() Model
}
