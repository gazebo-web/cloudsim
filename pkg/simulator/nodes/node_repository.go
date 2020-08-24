package nodes

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Repository
type Repository interface {
	Create(node Node) (*Node, error)
}

// repository
type repository struct {
	db *gorm.DB
}

// NewRepository
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Create
func (r *repository) Create(node Node) (*Node, error) {
	if pk := r.db.NewRecord(node); !pk {
		return nil, errors.New("Node already exists")
	}
	if err := r.db.Create(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}
