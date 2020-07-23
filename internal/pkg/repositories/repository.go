package repositories

import (
	"errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
)

var (
	ErrNoFilter = errors.New("no filters provided")
	ErrNoEntitiesUpdated = errors.New("no entities were updated")
	ErrNoEntitiesDeleted = errors.New("no entities were deleted")
)

// Repository represents a generic repository layer interface.
type Repository interface {
	Create([]domain.Entity) ([]domain.Entity, error)
	Find(output interface{}, offset, limit *int, filters ...Filter) error
	FindOne(entity domain.Entity, filters ...Filter) error
	Update(data interface{}, filters ...Filter) error
	Delete(filters ...Filter) error
	SingularName() string
	PluralName() string
	Model() domain.Entity
}
