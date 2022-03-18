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

// Create inserts a single entry.
//	entity: The entry to insert.
func (r *repositorySQL) Create(entity Model) (Model, error) {
	result, err := r.CreateBulk([]Model{entity})
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

// CreateBulk is a bulk operation to create multiple entries with a single operation.
//	entities: should be a slice of the same data structure implementing Model.
func (r *repositorySQL) CreateBulk(entities []Model) ([]Model, error) {
	for _, entity := range entities {
		err := r.DB.Model(r.Model()).Create(entity).Error
		if err != nil {
			return nil, err
		}
	}
	return entities, nil
}

// Find filters entries and stores filtered entries in output.
//	output: will contain the result of the query. It must be a pointer to a slice.
//	offset: defines the number of results to skip before loading values to output.
//	limit: defines the maximum number of entries to return. A nil value returns infinite results.
// 	filters: filter entries by field value.
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

// FindOne filters entries and stores the first filtered entry in output, it must be a pointer to
// a data structure implementing Model.
func (r *repositorySQL) FindOne(output Model, filters ...Filter) error {
	if len(filters) == 0 {
		return ErrNoFilter
	}
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.First(output)
	err := q.Error
	if err != nil {
		return err
	}
	return nil
}

// Update updates all model entries that match the provided filters with the given data.
//	data: must be a map[string]interface{}
//  filters: filter entries that should be updated.
func (r *repositorySQL) Update(data interface{}, filters ...Filter) error {
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.Update(data)
	err := q.Error
	if err != nil {
		return err
	} else if q.RowsAffected == 0 {
		return ErrNoEntriesUpdated
	}
	return nil
}

// Delete removes all the model entries that match filters.
//  filters: filter entries that should be deleted.
func (r *repositorySQL) Delete(filters ...Filter) error {
	q := r.startQuery()
	q = r.setQueryFilters(q, filters)
	q = q.Delete(r.Model())
	err := q.Error
	if err != nil {
		return err
	} else if q.RowsAffected == 0 {
		return ErrNoEntriesDeleted
	}
	return nil
}

// startQuery inits a gorm query for this repository's model. Multiple filters are ANDd together.
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
