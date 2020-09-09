package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// StartSimulationData has all the information needed to start a simulation.
// It's used as the data type for the action's store.
type StartSimulationData struct {
	GroupID               simulations.GroupID
	GazeboServerPodName   string
	CreateMachinesInputs  []cloud.CreateMachinesInput
	CreateMachinesOutputs []cloud.CreateMachinesOutput
	GazeboServerPodIP     string
	BaseLabels            map[string]string
	GazeboLabels          map[string]string
	BridgeLabels          map[string]string
	FieldComputerLabels   map[string]string
	GazeboNodeSelector    map[string]string
	GazeboPodResource     orchestrator.Resource
}
