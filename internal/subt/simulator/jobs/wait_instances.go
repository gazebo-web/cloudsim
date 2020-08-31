package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// WaitForInstances is used to wait until all required instances are ready.
var WaitForInstances = &actions.Job{
	Name:       "wait-for-instances",
	Execute:    waitForInstances,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// waitForInstances is the main process executed by WaitForInstances.
func waitForInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	// Parse group id.
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	// Create simulator context
	simCtx := context.NewContext(ctx)

	// Get machine list from context.
	ctxValue := simCtx.Value("machine-list")
	createMachinesOutput, ok := ctxValue.([]cloud.CreateMachinesOutput)

	// Create input
	var waitMachinesOkInputs []cloud.WaitMachinesOKInput
	for _, c := range createMachinesOutput {
		waitMachinesOkInputs = append(waitMachinesOkInputs, *c.ToWaitMachinesOKInput())
	}

	// Wait until machines are OK.
	err := simCtx.Platform().Machines().WaitOK(waitMachinesOkInputs)
	if err != nil {
		return nil, err
	}

	return gid, nil
}
