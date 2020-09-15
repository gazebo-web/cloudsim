package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// StartSimulationData contains all the information used across the start simulation action.
type StartSimulationData struct {
	GroupID simulations.GroupID
}
