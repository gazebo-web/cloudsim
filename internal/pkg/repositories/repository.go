package repositories

import (
	"errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
)

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

// Repository represents a generic repository layer interface.
type Repository interface {
	Create([]domain.Entity) ([]domain.Entity, error)
	Find(output interface{}, page, pageSize *int, filters ...Filter) error
	FindOne(entity domain.Entity, filters ...Filter) error
	Update(data interface{}, filters ...Filter) error
	Delete(filters ...Filter) error
	Count(filters ...Filter) (int, error)
	SingularName() string
	PluralName() string
	Model() domain.Entity
}
