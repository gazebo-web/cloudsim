package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type Repository struct {
	parent *simulations.Repository
}

type IRepository interface {
	CountByOwnerAndCircuit(owner, circuit string) (int, error)
	Get(groupID string) (*Simulation, error)
	Create(sim *simulations.Simulation) (*Simulation, error)
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	parent := simulations.NewRepository(db, "subt")
	repository := parent.(*simulations.Repository)
	r = &Repository{
		parent: repository,
	}
	return r
}

func (r *Repository) CountByOwnerAndCircuit(owner, circuit string) (int, error) {
	panic("Not implemented")
}

func (r *Repository) Get(groupID string) (*Simulation, error) {
	panic("Not implemented")
}

func (r *Repository) Create(sim *simulations.Simulation) (*Simulation, error) {
	subtSim := &Simulation{
		Base:                sim,
		GroupID:             sim.GroupID,
		Score:               nil,
		SimTimeDurationSec:  0,
		RealTimeDurationSec: 0,
		ModelCount:          0,
	}

	return subtSim, nil
}