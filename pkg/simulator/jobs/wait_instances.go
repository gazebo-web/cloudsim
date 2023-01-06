package jobs

import (
	"github.com/gazebo-web/cloudsim/pkg/actions"
	"github.com/gazebo-web/cloudsim/pkg/machines"
	"github.com/gazebo-web/cloudsim/pkg/simulator/state"
	"github.com/jinzhu/gorm"
)

// WaitForInstancesInput is the input of the WaitForInstances job.
type WaitForInstancesInput []machines.CreateMachinesOutput

// WaitForInstancesOutput is the output of the WaitForInstances job.
type WaitForInstancesOutput []machines.CreateMachinesOutput

// WaitForInstances is used to wait until all required instances are ready.
var WaitForInstances = &actions.Job{
	Name:    "wait-for-instances",
	Execute: waitForInstances,
}

// waitForInstances is the main process executed by WaitForInstances.
func waitForInstances(store actions.Store, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	// Parse input data
	machineList := value.(WaitForInstancesInput)

	// Create input
	var waitMachinesOkInputs []machines.WaitMachinesOKInput
	for _, c := range machineList {
		waitMachinesOkInputs = append(waitMachinesOkInputs, c.ToWaitMachinesOKInput())
	}

	// Create get platform namespace from state.
	s := store.State().(state.PlatformGetter)

	// Wait until machines are OK.
	err := s.Platform().Machines().WaitOK(waitMachinesOkInputs)
	if err != nil {
		return nil, err
	}

	return WaitForInstancesOutput(machineList), nil
}
