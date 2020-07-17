package repositories

import "github.com/jinzhu/gorm"

// gormRepository is a Repository implementation using GORM.
type gormRepository struct {
	db *gorm.DB
	model interface{}
}

// GetAll returns a slice of entities.
// If `offset` and `limit` are not nil, it will return up to `limit` results from the provided `offset`.
func (g gormRepository) GetAll(offset, limit *int) ([]interface{}, error) {
	var output []interface{}
	q := g.db.Model(&g.model)
	if offset != nil {
		q = q.Offset(*offset)
	}
	if limit != nil {
		q = q.Limit(*limit)
	}
	err := q.Find(&output).Error
	if err != nil {
		return nil, err
	}
	return output, nil
}

// Get returns a slice of entities with the given uuids.
func (g gormRepository) Get(uuids []string) ([]interface{}, error) {
	var out []interface{}
	q := g.db.Model(&g.model)
	for _, u := range uuids {
		q = q.Where("uuid = ?", u)
	}
	err := q.Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates the entities from the given uuids with the given data.
func (g gormRepository) Update(uuids []string, data interface{}) ([]interface{}, error) {
	q := g.db.Model(&g.model)
	for _, uuid := range uuids {
		q = q.Where("uuid = ?", uuid)
	}
	err := q.Update(&data).Error
	if err != nil {
		return nil, err
	}
	return []interface{}{uuids}, nil
}

// Delete removes entities with the given uuids.
func (g gormRepository) Delete(uuids []string) ([]interface{}, error) {
	q := g.db.Model(&g.model)
	for _, u := range uuids {
		q = q.Where("uuid = ?", u)
	}
	err := q.Delete(&g.model).Error
	if err != nil {
		return nil, err
	}
	return []interface{}{uuids}, nil
}

// NewSQLRepository initialez a new gormRepository.
func NewSQLRepository(db *gorm.DB, model interface{}) Repository {
	return &gormRepository{
		db: db,
		model: model,
	}
}