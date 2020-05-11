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

func (r *repository) Get(groupID string) (*simulations.Simulation, error) {
	panic("Not implemented")
}

func (r *repository) Create(simulation *simulations.Simulation) (*simulations.Simulation, error) {
	r.Repository.Create(simulation)

	sim := &Simulation{
		Base:                simulation,
		GroupID:             simulation.GroupID,
		Score:               nil,
		SimTimeDurationSec:  0,
		RealTimeDurationSec: 0,
		ModelCount:          0,
	}

	return subtSim, nil
}