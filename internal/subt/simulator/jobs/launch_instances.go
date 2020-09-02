package jobs

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
)

// LaunchInstances is used to launch the required instances to run a simulation.
var LaunchInstances = &actions.Job{
	Name:            "launch-instances",
	PreHooks:        []actions.JobFunc{createMachineInputs},
	Execute:         launchInstances,
	RollbackHandler: rollbackLaunchInstances,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// createMachineInputs creates the needed input for launchInstances.
func createMachineInputs(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	subnet, zone := simCtx.Platform().Store().Machines().SubnetAndZone()

	input := []cloud.CreateMachinesInput{
		{
			InstanceProfile: simCtx.Platform().Store().Machines().InstanceProfile(),
			KeyName:         simCtx.Platform().Store().Machines().KeyName(),
			Type:            simCtx.Platform().Store().Machines().Type(),
			Image:           simCtx.Platform().Store().Machines().BaseImage(),
			MinCount:        1,
			MaxCount:        1,
			FirewallRules:   simCtx.Platform().Store().Machines().FirewallRules(),
			SubnetID:        subnet,
			Zone:            zone,
			Tags:            simCtx.Platform().Store().Machines().Tags(sim, "gzserver", "gzserver"),
			InitScript:      simCtx.Platform().Store().Machines().InitScript(),
			Retries:         10,
		},
	}

	robots, err := simCtx.Services().Simulations().GetRobots(gid)
	for _, r := range robots {
		tags := simCtx.
			Platform().
			Store().
			Machines().
			Tags(sim, "field-computer", fmt.Sprintf("fc-%s", r.Name()))

		input = append(input, cloud.CreateMachinesInput{
			InstanceProfile: simCtx.Platform().Store().Machines().InstanceProfile(),
			KeyName:         simCtx.Platform().Store().Machines().KeyName(),
			Type:            simCtx.Platform().Store().Machines().Type(),
			Image:           simCtx.Platform().Store().Machines().BaseImage(),
			MinCount:        1,
			MaxCount:        1,
			FirewallRules:   simCtx.Platform().Store().Machines().FirewallRules(),
			SubnetID:        subnet,
			Zone:            zone,
			Tags:            tags,
			InitScript:      simCtx.Platform().Store().Machines().InitScript(),
			Retries:         10,
		})
	}

	return map[string]interface{}{
		"groupID":              gid,
		"simulation":           sim,
		"createMachinesInputs": input,
	}, nil
}

// launchInstances is the main process executed by LaunchInstances.
func launchInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	// Get map from prehook
	inputMap := value.(map[string]interface{})

	// Parse group id
	gid, ok := inputMap["groupID"].(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	// Parse create machines input
	createMachineInputs, ok := inputMap["createMachinesInputs"].([]cloud.CreateMachinesInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Set the minimum amount of machines that are required to launch the simulation
	var minRequested int
	for _, c := range createMachineInputs {
		minRequested += int(c.MinCount)
	}

	// Initialize sim context
	simCtx := context.NewContext(ctx)

	// Get the amount of reserved machines at the moment.
	reserved := simCtx.Platform().Machines().Count(cloud.CountMachinesInput{
		Filters: map[string][]string{
			"tag:cloudsim-simulation-worker": {
				simCtx.Platform().Store().Machines().NamePrefix(),
			},
			"instance-state-name": {
				"pending",
				"running",
			},
		},
	})

	// Create new wait request to check instances availability.
	req := waiter.NewWaitRequest(func() (bool, error) {
		reserved = simCtx.Platform().Machines().Count(cloud.CountMachinesInput{
			Filters: map[string][]string{
				"tag:cloudsim-simulation-worker": {
					simCtx.Platform().Store().Machines().NamePrefix(),
				},
				"instance-state-name": {
					"pending",
					"running",
				},
			},
		})
		if reserved == -1 {
			return false, errors.New("error waiting for instances")
		}
		return minRequested > simCtx.Platform().Store().Machines().Limit()-reserved, nil
	})

	// Get timeout and poll frequency from store
	timeout := simCtx.Platform().Store().Machines().Timeout()
	pollFreq := simCtx.Platform().Store().Machines().PollFrequency()

	// Execute request
	err := req.Wait(timeout, pollFreq)
	if err != nil {
		return nil, err
	}

	// Get simulation from input.
	sim, ok := inputMap["simulation"].(simulations.Simulation)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Create instances
	output, err := simCtx.Platform().Machines().Create(createMachineInputs)
	if len(output) > 0 {
		if dataErr := deployment.SetJobData(tx, nil, "machine-list", output); dataErr != nil {
			return nil, dataErr
		}
	}
	if err != nil {
		return nil, err
	}

	// Add machine list to context
	// TODO: Create constant for 'machine-list'
	simCtx = context.WithValue(simCtx, "machine-list", output)

	return gid, nil
}

// rollbackLaunchInstances is the process in charge of terminating the machine instances that were created in launchInstances.
func rollbackLaunchInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {
	jobData, dataErr := deployment.GetJobData(tx, &deployment.CurrentJob, "machine-list")
	if dataErr != nil {
		return nil, err
	}

	machineList, ok := jobData.([]cloud.CreateMachinesOutput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	simCtx := context.NewContext(ctx)
	for _, m := range machineList {
		if m.ToTerminateMachinesInput() != nil {
			_ = simCtx.Platform().Machines().Terminate(*m.ToTerminateMachinesInput())
		}
	}
	return nil, nil
}
