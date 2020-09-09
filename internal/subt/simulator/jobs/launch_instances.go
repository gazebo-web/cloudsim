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

const dataMachineListKey = "machine-list"

// LaunchInstances is used to launch the required instances to run a simulation.
var LaunchInstances = &actions.Job{
	Name:            "launch-instances",
	PreHooks:        []actions.JobFunc{createMachineInputs},
	Execute:         launchInstances,
	PostHooks:       []actions.JobFunc{launchInstancesPostHook},
	RollbackHandler: rollbackLaunchInstances,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// createMachineInputs creates the needed input for launchInstances.
func createMachineInputs(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	simCtx := context.NewContext(ctx)

	storeData := simCtx.Store().Get().(*StartSimulationData)

	sim, err := simCtx.Services().Simulations().Get(storeData.GroupID)
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

	robots, err := simCtx.Services().Simulations().GetRobots(storeData.GroupID)
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

	storeData.CreateMachinesInputs = input

	err = simCtx.Store().Set(storeData)
	if err != nil {
		return nil, err
	}

	return input, nil
}

// launchInstances is the main process executed by LaunchInstances.
func launchInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	// Get map from prehook
	input, ok := value.([]cloud.CreateMachinesInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Set the minimum amount of machines that are required to launch the simulation
	var minRequested int
	for _, c := range input {
		minRequested += int(c.MinCount)
	}

	// Initialize sim context
	simCtx := context.NewContext(ctx)

	// Initialize filters
	filters := map[string][]string{
		"tag:cloudsim-simulation-worker": {
			simCtx.Platform().Store().Machines().NamePrefix(),
		},
		"instance-state-name": {
			"pending",
			"running",
		},
	}

	// Get the amount of reserved machines at the moment.
	reserved := simCtx.Platform().Machines().Count(cloud.CountMachinesInput{
		Filters: filters,
	})

	if minRequested <= simCtx.Platform().Store().Machines().Limit()-reserved {
		err := waitForMinRequestedInstances(simCtx, filters, minRequested)
		if err != nil {
			return nil, err
		}
	}

	// Create instances
	output, err := simCtx.Platform().Machines().Create(input)

	return map[string]interface{}{
		"output": output,
		"error":  err,
	}, nil
}

func waitForMinRequestedInstances(ctx context.Context, filters map[string][]string, minRequested int) error {
	// Create new wait request to check instances availability.
	req := waiter.NewWaitRequest(func() (bool, error) {
		reserved := ctx.Platform().Machines().Count(cloud.CountMachinesInput{
			Filters: filters,
		})
		if reserved == -1 {
			return false, errors.New("error waiting for instances")
		}
		return minRequested > ctx.Platform().Store().Machines().Limit()-reserved, nil
	})

	// Get timeout and poll frequency from store
	timeout := ctx.Platform().Store().Machines().Timeout()
	pollFreq := ctx.Platform().Store().Machines().PollFrequency()

	// Execute request
	err := req.Wait(timeout, pollFreq)
	if err != nil {
		return err
	}
	return nil
}

func launchInstancesPostHook(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	data := ctx.Store().Get().(*StartSimulationData)

	execMap := value.(map[string]interface{})

	output := execMap["output"].([]cloud.CreateMachinesOutput)
	err := execMap["error"].(error)

	// Persist machine list if there are more than 0.
	if len(output) > 0 {
		if err := persistMachineList(ctx, data, output); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return data.GroupID, nil
}

func persistMachineList(ctx actions.Context, data *StartSimulationData, machineList []cloud.CreateMachinesOutput) error {
	data.CreateMachinesOutputs = machineList
	if err := ctx.Store().Set(data); err != nil {
		return err
	}
	return nil
}

// rollbackLaunchInstances is the process in charge of terminating the machine instances that were created in launchInstances.
func rollbackLaunchInstances(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	data := ctx.Store().Get().(*StartSimulationData)

	simCtx := context.NewContext(ctx)
	for _, output := range data.CreateMachinesOutputs {
		if output.ToTerminateMachinesInput() != nil {
			_ = simCtx.Platform().Machines().Terminate(*output.ToTerminateMachinesInput())
		}
	}
	return nil, nil
}
