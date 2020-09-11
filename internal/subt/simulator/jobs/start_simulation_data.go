package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// StartSimulationData contains all the information used across the start simulation action.
type StartSimulationData struct {
	GroupID            simulations.GroupID
	GazeboServerPod    orchestrator.Resource
	CreatedMachineList []cloud.CreateMachinesOutput
}
