package repositories

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
)

// Repository represents a generic repository layer interface.
type Repository interface {
	Create([]domain.Entity) ([]domain.Entity, error)
	Find(offset, limit *int, filters ...Filter) ([]domain.Entity, error)
	FindOne(filters ...Filter) (domain.Entity, error)
	Update(data domain.Entity, filters ...Filter) error
	Delete(filters ...Filter) error
	SingularName() string
	PluralName() string
	Model() domain.Entity
}
