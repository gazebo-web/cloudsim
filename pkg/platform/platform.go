package platform

import (
	"github.com/marcoshuck/cloudsim-refactor-proposal/pkg/platform/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
)

// Platform groups a set of components for creating simulations.
// Each application will make use of one platform to run their simulations.
// The cloudsim team provides a default Kubernetes and AWS implementation of this Platform.
// Other combinations could be implemented after adding their respective subcomponents.
type Platform interface {
	Storage() cloud.Storage
	Machines() cloud.Machines
	Orchestrator() orchestrator.Orchestrator
}
