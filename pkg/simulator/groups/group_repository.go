package groups

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type IRepository interface {
	Create(group Group) (*Group, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(group Group) (*Group, error) {
	if pk := r.db.NewRecord(group); !pk {
		return nil, errors.New("Group already exists")
	}
	if err := r.db.Create(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}
