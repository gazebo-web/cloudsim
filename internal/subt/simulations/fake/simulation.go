package fake

import (
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
)

// SimulationConfig is used to configure a fake subt fakeSimulation.
type SimulationConfig struct {
	GroupID simulations.GroupID
	Status  simulations.Status
	Kind    simulations.Kind
	Error   *simulations.Error
	Image   string
	Track   string
}

// fakeSimulation is a fake imulation implementation.
type fakeSimulation struct {
	simulations.Simulation
	track string
}

// Track returns the track of a simulation.
func (s *fakeSimulation) Track() string {
	return s.track
}

// NewSimulation initializes a new Simulation interface using a fake implementation.
func NewSimulation(config SimulationConfig) subt.Simulation {
	return &fakeSimulation{
		Simulation: fake.NewSimulation(config.GroupID, config.Status, config.Kind, config.Error, config.Image),
		track:      config.Track,
	}
}
