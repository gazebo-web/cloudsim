package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type repository struct {
	simulations.Repository
}

type Repository interface {
	simulations.Repository
	CountByOwnerAndCircuit(owner, circuit string) (int, error)
	Aggregate(simulation *Simulation) (*Simulation, error)
}

func NewRepository(db *gorm.DB) Repository {
	parent := simulations.NewRepository(db, "subt")
	r := parent.(simulations.Repository)
	return &repository{
		Repository: r,
	}
}

func (r *repository) CountByOwnerAndCircuit(owner, circuit string) (int, error) {
	panic("Not implemented")
}

func (r *repository) Aggregate(simulation *Simulation) (*Simulation, error) {
	err := r.GetDB().Create(simulation).Error
	if err != nil {
		return nil, err
	}
	return simulation, nil
}