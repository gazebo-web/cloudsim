package repositories

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
)

// Repository represents a generic repository layer interface.
type Repository interface {
	Creator
	Finder
	Updater
	Remover
	Namer
	Modeler
}

// Creator has a method to create a slice of entities.
type Creator interface {
	Create([]domain.Entity) ([]domain.Entity, error)
}

// Finder has methods to find an entity or a slice of entities.
type Finder interface {
	Find(output interface{}, offset, limit *int, filters ...Filter) error
	FindOne(entity domain.Entity, filters ...Filter) error
}

// Updater has a method to update entities.
type Updater interface {
	Update(data domain.Entity, filters ...Filter) error
}

// Remover has a method to delete entities.
type Remover interface {
	Delete(filters ...Filter) error
}

// Namer has methods to name an entity.
type Namer interface {
	SingularName() string
	PluralName() string
}

// Modeler has a method to return the repository's model.
type Modeler interface {
	Model() domain.Entity
}
