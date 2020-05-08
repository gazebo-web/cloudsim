package circuits

import (
	"github.com/jinzhu/gorm"
)

type IRepository interface {
	GetByName(name string) (*Circuit, error)
}

type Repository struct {
	Db        *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	r = &Repository{
		Db:        db,
	}
	return r
}

func (r *Repository) GetByName(name string) (*Circuit, error) {
	panic("implement me")
}