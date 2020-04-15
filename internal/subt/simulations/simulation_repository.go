package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type Repository struct {
	*simulations.Repository
}

type IRepository interface {
	simulations.IRepository
	CountByOwnerAndCircuit(owner, circuit string) (int, error)
}

func (r *Repository) CountByOwnerAndCircuit(owner, circuit string) (int, error) {
	panic("Not implemented")
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	parent := simulations.NewRepository(db)
	repository := parent.(*simulations.Repository)
	r = &Repository{
		Repository: repository,
	}
	return r
}
