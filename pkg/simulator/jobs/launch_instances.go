package jobs

import (
  "fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// LaunchInstancesInput is the input of the LaunchInstances job.
// It's used to pass the list of instances to create.
type LaunchInstancesInput []cloud.CreateMachinesInput

// LaunchInstancesOutput is the output of the LaunchInstances job.
// It's used to pass the list of instances created.
type LaunchInstancesOutput []cloud.CreateMachinesOutput

// LaunchInstances is a generic job to launch instances.
// It includes a rollback handler to terminate the instances that were created in this job.
var LaunchInstances = &actions.Job{
	Execute:         launchInstances,
	RollbackHandler: removeCreatedInstances,
}

// jobLaunchInstancesDataKey is the key used to persist the list of machines that were created in the LaunchInstances job.
const jobLaunchInstancesDataKey = "created-machines"

func launchInstances(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Get the store
	s := store.State().(state.PlatformGetter)

	// Parse the input
	in := value.(LaunchInstancesInput)

  fmt.Printf("Creating a machine...maybe??\n")
	// Trigger the machine creation.
	out, err := s.Platform().Machines().Create(in)

	// Set job data with the list of instances
	if dataErr := deployment.SetJobData(tx, nil, jobLaunchInstancesDataKey, LaunchInstancesOutput(out)); err != nil {
		return nil, dataErr
	}

	// Check for errors
	if err != nil {
		return nil, err
	}

	return LaunchInstancesOutput(out), nil
}

func removeCreatedInstances(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	// Get the store
	s := store.State().(state.PlatformGetter)

	// Get the list of instances from the execute function.
	data, dataErr := deployment.GetJobData(tx, nil, jobLaunchInstancesDataKey)
	if dataErr != nil {
		return nil, dataErr
	}

	// Parse the list of instances
	createdInstances := data.(LaunchInstancesOutput)

	// Terminate the instances
	for _, c := range createdInstances {
		_ = s.Platform().Machines().Terminate(c.ToTerminateMachinesInput())
	}

	return nil, nil
}
