package nodes

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// IRepository
type IRepository interface{
	Create(node Node) (*Node, error)
}

// Repository
type Repository struct{
	db *gorm.DB

}

// NewRepository
func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		db: db,
	}
}

// Create
func (r *Repository) Create(node Node) (*Node, error) {
	if pk := r.db.NewRecord(node); !pk {
		return nil, errors.New("Node already exists")
	}
	if err := r.db.Create(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}
