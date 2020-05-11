package simulator

import (
	"context"
	"errors"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/nodes"
	"sync"
	"time"
)

// Simulator
type Simulator interface {
	appendRunningSimulation(simulation *RunningSimulation)
	Recover(ctx context.Context, getApplicationLabel func() *string, getGazeboConfig func(sim *simulations.Simulation) GazeboConfig) error
	GetRunningSimulation(groupID string) *RunningSimulation
	GetRunningSimulations() map[string]*RunningSimulation
	SetRunningSimulations(simulations *map[string]*RunningSimulation) error
	RLock()
	RUnlock()
	Lock()
	Unlock()
}

// Config represents a set of options to configure a simulator.
type Config struct {
	NamePrefix               string `env:"AWS_INSTANCE_NAME_PREFIX,required"`
	ShouldTerminateInstances bool   `env:"EC2_NODE_MGR_TERMINATE_INSTANCES" envDefault:"true"`
	IamInstanceProfile       string `env:"AWS_IAM_INSTANCE_PROFILE_ARN" envDefault:"arn:aws:iam::200670743174:instance-profile/cloudsim-ec2-node"`
	JoinCmd                  string `env:"KUBEADM_JOIN,required"`
	AvailableEC2Machines     int    `env:"IGN_EC2_MACHINES_LIMIT" envDefault:"-1"`
}

// simulator is the component responsible of creating the nodes
// and registering them in the kubernetes master.
type simulator struct {
	orchestrator           orchestrator.Kubernetes
	cloud                  *cloud.AmazonWS
	runningSimulations     map[string]*RunningSimulation
	lockRunningSimulations sync.RWMutex
	config                 Config
	repositories           repositories
	services               services
	Controller             IController
}

// services
type services struct {
	simulations simulations.Service
	simulator   IService
}

// repositories
type repositories struct {
	node nodes.Repository
}

// NewSimulatorInput
type NewSimulatorInput struct {
	Orchestrator orchestrator.Kubernetes
	Cloud        *cloud.AmazonWS
	Db           *gorm.DB
}

// NewSimulator returns a new simulator instance.
func NewSimulator(input NewSimulatorInput) Simulator {
	cfg := Config{}
	if err := env.Parse(cfg); err != nil {
		// TODO: Throw an error. Logger? Log Fatal?
	}
	s := simulator{
		orchestrator: input.Orchestrator,
		cloud:        input.Cloud,
		repositories: repositories{
			node: nodes.NewRepository(input.Db),
		},
		config: cfg,
	}
	s.services.simulator = NewSimulatorService(s.repositories.node)
	return &s
}

// GetRunningSimulation
func (s *simulator) GetRunningSimulation(groupID string) *RunningSimulation {
	return s.runningSimulations[groupID]
}

// GetRunningSimulations
func (s *simulator) GetRunningSimulations() map[string]*RunningSimulation {
	return s.runningSimulations
}

// SetRunningSimulations
func (s *simulator) SetRunningSimulations(simulations *map[string]*RunningSimulation) error {
	if simulations == nil {
		return errors.New("SetRunningSimulations cannot receive a nil argument")
	}
	s.Lock()
	defer s.Unlock()
	s.runningSimulations = *simulations
	return nil
}

// appendRunningSimulation adds a new running simulation to the map of running simulations.
func (s *simulator) appendRunningSimulation(simulation *RunningSimulation) {
	s.Lock()
	defer s.Unlock()
	if s.runningSimulations[simulation.GroupID] != nil {
		return
	}
	s.runningSimulations[simulation.GroupID] = simulation
}

// RestoreRunningSimulation
func (s *simulator) RestoreRunningSimulation(ctx context.Context, simulation *simulations.Simulation, config GazeboConfig) error {
	validFor, err := time.ParseDuration(*simulation.ValidFor)
	if err != nil {
		return err
	}
	input := NewRunningSimulationInput{
		GroupID:          *simulation.GroupID,
		Owner:            *simulation.Owner,
		MaxSeconds:       config.MaxSeconds,
		ValidFor:         validFor,
		worldStatsTopic:  config.WorldStatsTopic,
		worldWarmupTopic: config.WorldWarmupTopic,
	}
	rs, err := NewRunningSimulation(ctx, input)
	if err != nil {
		return err
	}
	s.appendRunningSimulation(rs)
	return nil
}

// Recover
func (s *simulator) Recover(ctx context.Context, getApplicationLabel func() *string, getGazeboConfig func(sim *simulations.Simulation) GazeboConfig) error {
	label := getApplicationLabel()
	pods, err := s.orchestrator.GetAllPods(label)
	if err != nil {
		logger.Logger(ctx).Error("[SIMULATOR|RECOVER] Error getting initial list of pods from orchestrator", err)
		return err
	}

	runningSims := make(map[string]bool)
	for _, p := range pods {
		running, found := runningSims[p.GroupID]
		if !found {
			running = true
		}
		runningSims[p.GroupID] = running && p.IsRunning
	}

	for groupID, running := range runningSims {
		if !running {
			continue
		}

		sim, err := s.services.simulations.Get(groupID)
		if err != nil {
			return err
		}

		if !simulations.StatusRunning.Equal(*sim.Status) {
			continue
		}

		if err := s.RestoreRunningSimulation(ctx, sim, getGazeboConfig(sim)); err != nil {
			return err
		}
	}
	return nil
}

// RLock
func (s *simulator) RLock() {
	s.lockRunningSimulations.RLock()
}

// RUnlock
func (s *simulator) RUnlock() {
	s.lockRunningSimulations.RUnlock()
}

// Lock
func (s *simulator) Lock() {
	s.lockRunningSimulations.Lock()
}

// Unlock
func (s *simulator) Unlock() {
	s.lockRunningSimulations.Unlock()
}
