package fake

import (
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"time"
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
	name       string
	track      string
	token      *string
	robots     []simulations.Robot
	marsupials []simulations.Marsupial
}

// GetName returns the simulation name.
func (s *simulation) GetName() string {
	return s.name
}

// GetWorldIndex returns the world index.
func (s *simulation) GetWorldIndex() int {
	return 0
}

// GetToken returns the access token of a simulation.
func (s *simulation) GetToken() *string {
	return s.token
}

// GetRobots return the list of simulations.Robot that will run in the simulation.
func (s *simulation) GetRobots() []simulations.Robot {
	return s.robots
}

// GetMarsupials return the list of simulations.Marsupial for the simulation.
func (s *simulation) GetMarsupials() []simulations.Marsupial {
	return s.marsupials
}

// GetTrack returns the track of a simulation.
func (s *simulation) GetTrack() string {
	return s.track
}

// NewSimulation initializes a new Simulation interface using a fake implementation.
func NewSimulation(config SimulationConfig) subt.Simulation {
	return &simulation{
		Simulation: fake.NewSimulation(config.GroupID, config.Status, config.Kind, config.Error, config.Image, 1*time.Minute),
		track:      config.Track,
	}
}
