package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"time"
)

// Simulation represents the simulation that will be launched in the cloud.
// A copy of this entity could be found in the following path:
// simulations/models.go:15
type Simulation struct {
	gorm.Model

	// Add simulation fields here
}

func (s *Simulation) TableName() string {
	return "simulations"
}

func (s *Simulation) SingularName() string {
	return "simulation"
}

func (s *Simulation) PluralName() string {
	return "simulations"
}

func (s *Simulation) GetGroupID() simulations.GroupID {
	panic("implement me")
}

func (s *Simulation) GetStatus() simulations.Status {
	panic("implement me")
}

func (s *Simulation) HasStatus(status simulations.Status) bool {
	panic("implement me")
}

func (s *Simulation) SetStatus(status simulations.Status) {
	panic("implement me")
}

func (s *Simulation) GetKind() simulations.Kind {
	panic("implement me")
}

func (s *Simulation) IsKind(kind simulations.Kind) bool {
	panic("implement me")
}

func (s *Simulation) GetError() *simulations.Error {
	panic("implement me")
}

func (s *Simulation) GetImage() string {
	panic("implement me")
}

func (s *Simulation) GetValidFor() time.Duration {
	panic("implement me")
}

func NewSimulation() simulations.Simulation {
	return &Simulation{}
}
