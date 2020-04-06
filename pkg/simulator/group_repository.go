package simulator

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type IGroupRepository interface {
	Crete(group Group) (*Group, error)
}

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) Create(group Group) (*Group, error) {
	if pk := r.db.NewRecord(group); !pk {
		return nil, errors.New("Group already exists")
	}
	if err := r.db.Create(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}
