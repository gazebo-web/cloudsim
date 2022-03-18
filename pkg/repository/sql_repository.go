package repository

import (
	"github.com/jinzhu/gorm"
	"reflect"
)

// NewRepositorySQL initializes a new Repository implementation for SQL databases.
func NewRepositorySQL(db *gorm.DB, entity Model) Repository {
	return &repositorySQL{
		DB:     db,
		entity: entity,
	}
}

// repositorySQL implements Repository using gorm to support SQL databases.
type repositorySQL struct {
	DB     *gorm.DB
	entity Model
}

// Create is a bulk operation to create N amount of models passed in an array.
func (r *repositorySQL) Create(entities []Model) ([]Model, error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, entity := range entities {
		err := tx.Model(r.Model()).Create(entity).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// Find writes on output the result of applying offset, limit and filters to find a range of models.
func (r *repositorySQL) Find(output interface{}, offset, limit *int, filters ...Filter) error {
	q := r.startQuery()
	if limit != nil {
		q = q.Limit(*limit)
		if offset != nil {
			q = q.Offset(*offset)
		}
	}

	q = r.setQueryFilters(q, filters)
	q = q.Find(output)
	err := q.Error
	if err != nil {
		return err
	}
	return nil
}

// FindOne applies the given filters to write on entity the first result found.
func (r *repositorySQL) FindOne(entity Model, filters ...Filter) error {
	if len(filters) == 0 {
		return ErrNoFilter
	}
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.First(entity)
	err := q.Error
	if err != nil {
		return err
	}
	return nil
}

// Update updates with the given data the different models that match filters.
func (r *repositorySQL) Update(data interface{}, filters ...Filter) error {
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.Update(data)
	err := q.Error
	if err != nil {
		return err
	} else if q.RowsAffected == 0 {
		return ErrNoEntitiesUpdated
	}
	return nil
}

// Delete removes all the models that match filters.
func (r *repositorySQL) Delete(filters ...Filter) error {
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.Delete(r.Model())
	err := q.Error
	if err != nil {
		return err
	} else if q.RowsAffected == 0 {
		return ErrNoEntitiesDeleted
	}
	return nil
}

// startQuery inits a gorm query for this repository's model.
func (r *repositorySQL) startQuery() *gorm.DB {
	return r.DB.Model(r.Model())
}

// setQueryFilters applies the given filters to a gorm query.
func (r *repositorySQL) setQueryFilters(q *gorm.DB, filters []Filter) *gorm.DB {
	for _, f := range filters {
		q = q.Where(f.Template, f.Values...)
	}
	return q
}

// Model returns this repository's model.
func (r *repositorySQL) Model() Model {
	entity := reflect.ValueOf(r.entity)

	if entity.Kind() == reflect.Ptr {
		entity = reflect.Indirect(entity)
	}

	return reflect.New(entity.Type()).Interface().(Model)
}
