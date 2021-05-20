package gorm

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"reflect"
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

// repository is a Repository implementation using GORM.
type repository struct {
	DB     *gorm.DB
	Logger ign.Logger
	entity domain.Entity
}

// SingularName returns the singular name for this repository's entity.
// Example: "Car"
func (r repository) SingularName() string {
	return r.Model().SingularName()
}

// PluralName returns the plural name for this repository's entity.
// Example: "Cars"
func (r repository) PluralName() string {
	return r.Model().PluralName()
}

// Model returns a pointer to the entity struct for this repository.
// Example: &Car{}
func (r repository) Model() domain.Entity {
	entity := reflect.ValueOf(r.entity)

	if entity.Kind() == reflect.Ptr {
		entity = reflect.Indirect(entity)
	}

	return reflect.New(entity.Type()).Interface().(domain.Entity)
}

// Create persists the given entities in the database.
func (r repository) Create(entities []domain.Entity) ([]domain.Entity, error) {
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Creating %s. Input: %+v",
		r.SingularName(), r.PluralName(), entities))
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, entity := range entities {
		err := tx.Model(r.Model()).Create(entity).Error
		if err != nil {
			r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Creating %s failed. Input: %+v. Rolling back...",
				r.SingularName(), r.PluralName(), entity))
			tx.Rollback()
			return nil, err
		}
	}
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Creating %s succeeded. Committing...",
		r.SingularName(), r.PluralName()))
	err := tx.Commit().Error
	if err != nil {
		r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Error while commiting the transaction.",
			r.SingularName()))
		return nil, err
	}
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] %s created.",
		r.SingularName(), r.PluralName()))
	return entities, nil
}

// Find returns a list of entities that match the given filters.
// If `offset` and `limit` are not nil, it will return up to `limit` results from the provided `offset`.
func (r repository) Find(output interface{}, limit, offset *int, filters ...domain.Filter) error {
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting all %s. Filters: %+v",
		r.SingularName(), r.PluralName(), filters))
	q := r.startQuery()
	if limit != nil {
		r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Limit: %d.",
			r.SingularName(), *limit))
		q = q.Limit(*limit)
		if offset != nil {
			r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Offset: %d.",
				r.SingularName(), *offset))
			q = q.Offset(*offset)
		}
	}

	q = r.setQueryFilters(q, filters)
	q = q.Find(output)
	err := q.Error
	if err != nil {
		r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting all %s failed. Error: %+v",
			r.SingularName(), r.PluralName(), err))
		return err
	}
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting all %s succeed. Output: %+v",
		r.SingularName(), r.PluralName(), output))
	return nil
}

// FindOne returns an Entity that matches the given Filter.
func (r repository) FindOne(entity domain.Entity, filters ...domain.Filter) error {
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting %s. Filters: %+v",
		r.SingularName(), r.SingularName(), filters))
	if len(filters) == 0 {
		r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting %s failed. No filters provided.",
			r.SingularName(), r.SingularName()))
		return ErrNoFilter
	}
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.First(entity)
	err := q.Error
	if err != nil {
		r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting %s failed. Error: %+v.",
			r.SingularName(), r.SingularName(), err))
		return err
	}
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] %s found. Result: %+v.",
		r.SingularName(), r.SingularName(), entity))
	return nil
}

// Update updates the entities that match the given filters with the given data.
func (r repository) Update(data interface{}, filters ...domain.Filter) error {
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Updating with data: %+v. Filters: %+v",
		r.SingularName(), data, filters))
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.Update(data)
	err := q.Error
	if err != nil {
		r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Updating failed. Error: %+v",
			r.SingularName(), err))
		return err
	} else if q.RowsAffected == 0 {
		return ErrNoEntitiesUpdated
	}
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Updating succeed. Updated records: %d.",
		r.SingularName(), q.RowsAffected))
	return nil
}

// Delete removes a set of entities that match the given filters.
func (r repository) Delete(filters ...domain.Filter) error {
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Deleting records. Filters: %+v",
		r.SingularName(), filters))
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.Delete(r.Model())
	err := q.Error
	if err != nil {
		r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Deleting failed. Error: %+v",
			r.SingularName(), err))
		return err
	} else if q.RowsAffected == 0 {
		return ErrNoEntitiesDeleted
	}
	r.Logger.Debug(fmt.Sprintf(" [%s.Repository] Deleting succeed. Removed records: %d.",
		r.SingularName(), q.RowsAffected))
	return nil
}

func (r repository) startQuery() *gorm.DB {
	return r.DB.Model(r.Model())
}

func (r repository) setQueryFilters(q *gorm.DB, filters []domain.Filter) *gorm.DB {
	for _, f := range filters {
		q = q.Where(f.Template(), f.Values()...)
	}
	return q
}

// NewRepository initializes a new Repository implementation using gorm.
func NewRepository(db *gorm.DB, logger ign.Logger, entity domain.Entity) domain.Repository {
	return &repository{
		DB:     db,
		Logger: logger,
		entity: entity,
	}
}
