package circuits

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type IRepository interface {
	GetByName(name string) (*Circuit, error)
}

type Repository struct {
	Db        *gorm.DB
	whitelist map[string]bool
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	r = &Repository{
		Db:        db,
		whitelist: generateWhitelist(),
	}
	return r
}

func (r *Repository) GetByName(name string) (*Circuit, error) {
	panic("implement me")
}

func (r *Repository) GetFromWhitelist(name string) (*string, *bool, error) {
	allowed, ok := r.whitelist[name]
	if !ok {
		return nil, nil, errors.New("circuit doesn't exist")
	}
	return &name, &allowed, nil
}
