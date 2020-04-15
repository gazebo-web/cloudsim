package nodes

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type IRepository interface{
	Create(node Node) (*Node, error)
}

type Repository struct{
	db *gorm.DB

}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(node Node) (*Node, error) {
	if pk := r.db.NewRecord(node); !pk {
		return nil, errors.New("Node already exists")
	}
	if err := r.db.Create(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}
