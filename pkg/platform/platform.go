package platform

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
)

// Platform groups a set of components for creating simulations.
// Each application will make use of one platform to run their simulations.
// The cloudsim team provides a default Kubernetes and AWS implementation of this Platform.
// Other combinations could be implemented after adding their respective subcomponents.
type Platform interface {
	// Storage returns a cloud.Storage component.
	Storage() cloud.Storage

	// Machines returns a cloud.Machines component.
	Machines() cloud.Machines

	// Orchestrator returns a orchestrator.Cluster component.
	Orchestrator() orchestrator.Cluster

	// Store returns a store.Store component.
	Store() store.Store
}

// NewPlatform initializes a new platform using the given components.
func NewPlatform(machines cloud.Machines, storage cloud.Storage, cluster orchestrator.Cluster, store store.Store) Platform {
	return &platform{
		storage:      storage,
		machines:     machines,
		orchestrator: cluster,
		store:        store,
	}
}

// platform is a Platform implementation.
type platform struct {
	storage      cloud.Storage
	machines     cloud.Machines
	orchestrator orchestrator.Cluster
	store        store.Store
}

// Store returns a store.Store implementation.
func (p *platform) Store() store.Store {
	return p.store
}

// Storage returns a cloud.Storage implementation.
func (p *platform) Storage() cloud.Storage {
	return p.storage
}

// Machines returns a cloud.Machines implementation.
func (p *platform) Machines() cloud.Machines {
	return p.machines
}

// Orchestrator returns an orchestrator.Cluster implementation.
func (p *platform) Orchestrator() orchestrator.Cluster {
	return p.orchestrator
}
