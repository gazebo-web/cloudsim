package circuits

import (
	"github.com/jinzhu/gorm"
	"time"
)

type IRepository interface {
	GetByName(name string) (*Circuit, error)
	GetPending() ([]Circuit, error)
}

type Repository struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	r = &Repository{
		Db: db,
	}
	return r
}

func (r *Repository) GetByName(name string) (*Circuit, error) {
	var c Circuit
	err := r.Db.Model(&Circuit{}).Where("circuit = ?", name).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetPending() ([]Circuit, error) {
	var cs []Circuit
	err := r.Db.Model(&Circuit{}).Where("competition_date >= ?", time.Now()).Find(&cs).Error
	if err != nil {
		return nil, err
	}
	return cs, nil
}
