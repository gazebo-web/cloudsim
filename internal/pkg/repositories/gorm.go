package repositories

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"reflect"
)

// GormRepository is a Repository implementation using GORM.
type GormRepository struct {
	DB     *gorm.DB
	Logger ign.Logger
	entity domain.Entity
}

var (
	// ErrNegativePageSize is returned when a negative page size is passed to validatePagination
	ErrNegativePageSize = errors.New("negative page size")
	// ErrNegativePage is returned when a negative page is passed to validatePagination
	ErrNegativePage = errors.New("negative page")
)

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
	entity := reflect.ValueOf(g.entity)

	if entity.Kind() == reflect.Ptr {
		entity = reflect.Indirect(entity)
	}

	return reflect.New(entity.Type()).Interface().(domain.Entity)
}

// Create persists the given entities in the database.
func (g GormRepository) Create(entities []domain.Entity) ([]domain.Entity, error) {
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Creating %s. Input: %+v",
		g.SingularName(), g.PluralName(), entities))
	tx := g.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, entity := range entities {
		err := tx.Model(g.Model()).Create(entity).Error
		if err != nil {
			g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Creating %s failed. Input: %+v. Rolling back...",
				g.SingularName(), g.PluralName(), entity))
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
// If `page` and `pageSize` are not nil, it will return up to `pageSize` results from the provided `page`.
// NOTE: `page` number starts at 0.
func (g GormRepository) Find(output interface{}, page, pageSize *int, filters ...Filter) error {
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting all %s. Filters: %+v",
		g.SingularName(), g.PluralName(), filters))
	var err error
	page, pageSize, err = g.validatePagination(page, pageSize)
	if err != nil {
		return err
	}
	q := g.startQuery()
	limit, offset := g.calculatePagination(page, pageSize)
	if limit != nil {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Limit: %d.",
			g.SingularName(), *limit))
		q = q.Limit(*limit)
		if offset != nil {
			g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Offset: %d.",
				g.SingularName(), *offset))
			q = q.Offset(*offset)
		}
	}

	q = g.setQueryFilters(q, filters)
	q = q.Find(output)
	err = q.Error
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
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting %s. Filters: %+v",
		g.SingularName(), g.SingularName(), filters))
	if len(filters) == 0 {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting %s failed. No filters provided.",
			g.SingularName(), g.SingularName()))
		return ErrNoFilter
	}
	q := g.startQuery()
	q = g.setQueryFilters(q, filters)
	q = q.First(entity)
	err := q.Error
	if err != nil {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Getting %s failed. Error: %+v.",
			g.SingularName(), g.SingularName(), err))
		return err
	}
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] %s found. Result: %+v.",
		g.SingularName(), g.SingularName(), entity))
	return nil
}

// Update updates the entities that match the given filters with the given data.
func (g GormRepository) Update(data interface{}, filters ...Filter) error {
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Updating with data: %+v. Filters: %+v",
		g.SingularName(), data, filters))
	q := g.startQuery()
	q = g.setQueryFilters(q, filters)
	q = q.Update(data)
	err := q.Error
	if err != nil {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Updating failed. Error: %+v",
			g.SingularName(), err))
		return err
	} else if q.RowsAffected == 0 {
		return ErrNoEntitiesUpdated
	}
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Updating succeed. Updated records: %d.",
		g.SingularName(), q.RowsAffected))
	return nil
}

// Delete removes a set of entities that match the given filters.
func (g GormRepository) Delete(filters ...Filter) error {
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Deleting records. Filters: %+v",
		g.SingularName(), filters))
	q := g.startQuery()
	q = g.setQueryFilters(q, filters)
	q = q.Delete(g.Model())
	err := q.Error
	if err != nil {
		g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Deleting failed. Error: %+v",
			g.SingularName(), err))
		return err
	} else if q.RowsAffected == 0 {
		return ErrNoEntitiesDeleted
	}
	g.Logger.Debug(fmt.Sprintf(" [%s.Repository] Deleting succeed. Removed records: %d.",
		g.SingularName(), q.RowsAffected))
	return nil
}

// Count counts the amount of entities that match the given filters.
func (g GormRepository) Count(filters ...Filter) (int, error) {
	var count int
	q := g.startQuery()
	q = g.setQueryFilters(q, filters)
	err := q.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (g GormRepository) startQuery() *gorm.DB {
	return g.DB.Model(g.Model())
}

func (g GormRepository) setQueryFilters(q *gorm.DB, filters []Filter) *gorm.DB {
	for _, f := range filters {
		q = q.Where(f.Template(), f.Values()...)
	}
	return q
}

// calculatePagination returns the calculated SQL LIMIT and OFFSET from the given page and pageSize.
// NOTE: The count of `page` starts at 0.
func (g GormRepository) calculatePagination(page, pageSize *int) (limit *int, offset *int) {
	if pageSize == nil {
		return
	}

	// Set the limit if there is a pageSize
	limit = pageSize

	// Set the offset if a page number was defined
	var value int
	if page != nil {
		value = *page * *pageSize
		offset = &value
	}

	return
}

// validatePagination performs validation on `page` and `pageSize`.
// If `page` and `pageSize` are nil, it will assign the default values.
// page = 0
// pageSize = 10
func (g GormRepository) validatePagination(page, pageSize *int) (*int, *int, error) {
	var err error
	defaultPageSize := 10
	pageSize, err = g.checkPositivePaginationValue(pageSize, &defaultPageSize, ErrNegativePageSize)
	if err != nil {
		return nil, nil, err
	}

	defaultPage := 0
	page, err = g.checkPositivePaginationValue(page, &defaultPage, ErrNegativePage)
	if err != nil {
		return nil, nil, err
	}

	return page, pageSize, nil
}

// checkPositivePaginationValue checks if the given `value` is positive and not nil.
// If `value` is negative, returns the err passed as an argument.
// If `value` is nil, returns the default value passed as an argument.
// Otherwise, it returns the actual value.
func (g GormRepository) checkPositivePaginationValue(value *int, defaultValue *int, err error) (*int, error) {
	if value != nil {
		if *value < 0 {
			return nil, err
		}
		return value, nil
	}
	return defaultValue, nil
}


// NewGormRepository initializes a new Repository implementation using gorm.
func NewGormRepository(db *gorm.DB, logger ign.Logger, entity domain.Entity) Repository {
	return &GormRepository{
		DB:     db,
		Logger: logger,
		entity: entity,
	}
}
