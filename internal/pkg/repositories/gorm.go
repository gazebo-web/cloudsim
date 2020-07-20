package repositories

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// GormRepository is a Repository implementation using GORM.
type GormRepository struct {
	DB     *gorm.DB
	Logger ign.Logger
	entity domain.Entity
}

// SingularName returns the singular name for this repository's entity.
// Example: "Car"
func (g GormRepository) SingularName() string {
	return g.Model().SingularName()
}

// PluralName returns the plural name for this repository's entity.
// Example: "Cars"
func (g GormRepository) PluralName() string {
	return g.Model().PluralName()
}

// Model returns a pointer to the entity struct for this repository.
// Example: &Car{}
func (g GormRepository) Model() domain.Entity {
	c := new(domain.Entity)
	*c = g.entity
	return *c
}

// Create persists the given entities in the database.
func (g GormRepository) Create(entities []domain.Entity) ([]domain.Entity, error) {
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Creating %s. Input: %+v",
		g.SingularName(), g.PluralName(), entities))
	tx := g.DB.Begin()
	for _, e := range entities {
		err := tx.Model(g.Model()).Create(e).Error
		if err != nil {
			g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Creating %s failed. Input: %+v. Rolling back...",
				g.SingularName(), g.PluralName(), e))
			tx.Rollback()
			return nil, err
		}
	}
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Creating %s succeeded. Committing...",
		g.SingularName(), g.PluralName()))
	err := tx.Commit().Error
	if err != nil {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Error while commiting the transaction.",
			g.SingularName()))
		return nil, err
	}
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] %s created.",
		g.SingularName(), g.PluralName()))
	return entities, nil
}

// Find returns a list of entities that match the given filters.
// If `offset` and `limit` are not nil, it will return up to `limit` results from the provided `offset`.
func (g GormRepository) Find(output interface{}, offset, limit *int, filters ...Filter) error {
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting all %s. Filters: %+v",
		g.SingularName(), g.PluralName(), filters))
	q := g.startQuery()
	if offset != nil {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Offset: %d.",
			g.SingularName(), *offset))
		q = q.Offset(*offset)
	}
	if limit != nil {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Limit: %d.",
			g.SingularName(), *limit))
		q = q.Limit(*limit)
	}
	q = g.setQueryFilters(q, filters)
	err := q.Find(output).Error
	if err != nil {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting all %s failed. Error: %+v",
			g.SingularName(), g.PluralName(), err))
		return err
	}
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting all %s succeed. Output: %+v",
		g.SingularName(), g.PluralName(), output))
	return nil
}

// FindOne returns an Entity that matches the given Filter.
func (g GormRepository) FindOne(entity domain.Entity, filters ...Filter) error {
	if len(filters) == 0 {
		return errors.New("no filters provided")
	}
	q := g.startQuery()
	q = g.setQueryFilters(q, filters)
	err := q.First(entity).Error
	if err != nil {
		return err
	}
	return nil
}

// Update updates the entities that match the given filters with the given data.
func (g GormRepository) Update(data domain.Entity, filters ...Filter) error {
	q := g.startQuery()
	q = g.setQueryFilters(q, filters)
	err := q.Update(data).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a set of entities that match the given filters.
func (g GormRepository) Delete(filters ...Filter) error {
	q := g.startQuery()
	g.setQueryFilters(q, filters)
	err := q.Delete(g.Model()).Error
	if err != nil {
		return err
	}
	return nil
}

func (g GormRepository) startQuery() *gorm.DB {
	return g.DB.Model(g.Model())
}

func (g GormRepository) setQueryFilters(q *gorm.DB, filters []Filter) *gorm.DB {
	for _, f := range filters {
		q = q.Where(f.Template(), f.Value())
	}
	return q
}

// NewGormRepository initializes a new Repository implementation using gorm.
func NewGormRepository(db *gorm.DB, logger ign.Logger, entity domain.Entity) Repository {
	return &GormRepository{
		DB:     db,
		Logger: logger,
		entity: entity,
	}
}
