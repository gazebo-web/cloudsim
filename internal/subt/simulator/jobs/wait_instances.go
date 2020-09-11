package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// WaitForInstances is used to wait until all required instances are ready.
var WaitForInstances = &actions.Job{
	Name:       "wait-for-instances",
	Execute:    waitForInstances,
	InputType:  actions.GetJobDataType(&StartSimulationData{}),
	OutputType: actions.GetJobDataType(&StartSimulationData{}),
}

// waitForInstances is the main process executed by WaitForInstances.
func waitForInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	// Parse input data
	data := value.(*StartSimulationData)

	// Create input
	var waitMachinesOkInputs []cloud.WaitMachinesOKInput
	for _, c := range data.CreatedMachineList {
		waitMachinesOkInputs = append(waitMachinesOkInputs, *c.ToWaitMachinesOKInput())
	}

	// Create simulator context
	simCtx := context.NewContext(ctx)

	// Wait until machines are OK.
	err := simCtx.Platform().Machines().WaitOK(waitMachinesOkInputs)
	if err != nil {
		return nil, err
	}

	return data, nil
}
