package platform

import (
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/runsim"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
)

var (
	// ErrInvalidPlatformName is returned when a platform with an invalid name is created.
	ErrInvalidPlatformName = errors.New("invalid platform name")
)

// Platform groups a set of components for managing simulations.
// Each application will make use of one platform to run their simulations.
// The cloudsim team provides a default Kubernetes and AWS implementation of this Platform.
// Other combinations could be implemented after adding their respective subcomponents.
type Platform interface {
	// GetName returns the platform name.
	GetName() string

	// Storage returns a storage.Storage component.
	Storage() storage.Storage

	// Machines returns a machines.Machines component.
	Machines() machines.Machines

	// Orchestrator returns a orchestrator.Cluster component.
	Orchestrator() orchestrator.Cluster

	// Store returns a store.Store component.
	Store() store.Store

	// Secrets returns a secrets.Secrets component.
	Secrets() secrets.Secrets

	// RunningSimulations returns a runsim.Manager component.
	RunningSimulations() runsim.Manager

	// EmailSender returns an email.Sender component.
	EmailSender() email.Sender
}

// Components lists the components used to initialize a Platform.
type Components struct {
	Machines           machines.Machines
	Storage            storage.Storage
	Cluster            orchestrator.Cluster
	Store              store.Store
	Secrets            secrets.Secrets
	EmailSender        email.Sender
	RunningSimulations runsim.Manager
}

// NewPlatform initializes a new platform using the given components.
func NewPlatform(name string, components Components) (Platform, error) {
	if name == "" {
		return nil, ErrInvalidPlatformName
	}

	return &platform{
		name:               name,
		storage:            components.Storage,
		machines:           components.Machines,
		orchestrator:       components.Cluster,
		store:              components.Store,
		secrets:            components.Secrets,
		email:              components.EmailSender,
		runningSimulations: components.RunningSimulations,
	}, nil
}

// platform is a Platform implementation.
type platform struct {
	name               string
	storage            storage.Storage
	machines           machines.Machines
	orchestrator       orchestrator.Cluster
	store              store.Store
	secrets            secrets.Secrets
	runningSimulations runsim.Manager
	email              email.Sender
}

func (p *platform) GetName() string {
	return p.name
}

// EmailSender returns a email.Sender implementation.
func (p *platform) EmailSender() email.Sender {
	return p.email
}

// RunningSimulations returns a runsim.Manager implementation.
func (p *platform) RunningSimulations() runsim.Manager {
	return p.runningSimulations
}

// Store returns a store.Store implementation.
func (p *platform) Store() store.Store {
	return p.store
}

// Storage returns a storage.Storage implementation.
func (p *platform) Storage() storage.Storage {
	return p.storage
}

// Machines returns a machines.Machines implementation.
func (p *platform) Machines() machines.Machines {
	return p.machines
}

// Orchestrator returns an orchestrator.Cluster implementation.
func (p *platform) Orchestrator() orchestrator.Cluster {
	return p.orchestrator
}

// Secrets returns an secrets.Secrets implementation.
func (p *platform) Secrets() secrets.Secrets {
	return p.secrets
}
