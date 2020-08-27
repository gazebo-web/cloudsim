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
