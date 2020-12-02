package fake

import (
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
)

// SimulationConfig is used to configure a fake subt simulation.
type SimulationConfig struct {
	GroupID    simulations.GroupID
	Status     simulations.Status
	Kind       simulations.Kind
	Error      *simulations.Error
	Image      string
	Track      string
	Token      *string
	Robots     []simulations.Robot
	Marsupials []simulations.Robot
}

// simulation is a fake simulation implementation.
type simulation struct {
	simulations.Simulation
	track      string
	token      *string
	robots     []simulations.Robot
	marsupials []simulations.Marsupial
}

// Token returns the access token of a simulation.
func (s *simulation) GetToken() *string {
	return s.token
}

// Robots return the list of simulations.Robot that will run in the simulation.
func (s *simulation) GetRobots() []simulations.Robot {
	return s.robots
}

// Marsupials return the list of simulations.Marsupial for the simulation.
func (s *simulation) GetMarsupials() []simulations.Marsupial {
	return s.marsupials
}

// Track returns the track of a simulation.
func (s *simulation) GetTrack() string {
	return s.track
}

// NewSimulation initializes a new Simulation interface using a fake implementation.
func NewSimulation(config SimulationConfig) subt.Simulation {
	return &simulation{
		Simulation: fake.NewSimulation(config.GroupID, config.Status, config.Kind, config.Error, config.Image),
		track:      config.Track,
	}
}
