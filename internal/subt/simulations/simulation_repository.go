package simulations

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
)

type repository struct {
	simulations.Repository
}

type Repository interface {
	simulations.Repository
	CountByOwnerAndCircuit(owner, circuit string) (int, error)
	Aggregate(simulation *Simulation) (*Simulation, error)
}

func NewRepository(db *gorm.DB, platform string) Repository {
	parent := simulations.NewRepository(db, &platform, tools.Sptr("subt"))
	r := parent.(simulations.Repository)
	return &repository{
		Repository: r,
	}
}

func (r *repository) Create(simulation simulations.RepositoryCreateInput) (simulations.ServiceCreateOutput, error) {
	subtInput, ok := simulation.(RepositoryCreateInput)
	if !ok {
		return nil, errors.New("error casting")
	}

	sim := subtInput.Input()

	output, err := r.Repository.Create(sim)
	if err != nil {
		return nil, err
	}

	subtSim := subtInput.ChildInput()

	subtSim.GroupID = sim.GroupID
	subtSim.Base = sim

	_, err = r.Aggregate(subtSim)
	if err != nil {
		return nil, err
	}
	return output, nil
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
