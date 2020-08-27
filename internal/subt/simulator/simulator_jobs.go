package simulator

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"time"
)

// JobsStartSimulation groups the jobs needed to start a simulation.
var JobsStartSimulation = actions.Jobs{
	JobCheckPendingStatus,
	JobCheckSimulationParenthood,
	JobCheckParentSimulationWithError,
	JobUpdateSimulationStatusToLaunchNodes,
	JobLaunchNodes,
	JobWaitForNodes,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{}

// JobsRestartSimulation groups the jobs needed to restart a simulation.
var JobsRestartSimulation = actions.Jobs{}

//----------------------------------------------------------------------------------------------------------------------

// JobCheckPendingStatus is used to check that a certain simulation has pending status.
var JobCheckPendingStatus = &actions.Job{
	Name:       "check-pending-status",
	Execute:    checkPendingStatus,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// checkPendingStatus is the main process executed by JobCheckSimulationParenthood.
func checkPendingStatus(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := NewContext(ctx)

	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	if sim.Status() != simulations.StatusPending {
		return nil, simulations.ErrIncorrectStatus
	}
	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------

// JobCheckSimulationParenthood is used to check that a simulation is not of kind simulations.SimParent.
var JobCheckSimulationParenthood = &actions.Job{
	Name:       "check-simulation-parenthood",
	Execute:    checkSimulationIsParent,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// checkSimulationIsParent is the main process executed by JobCheckPendingStatus.
func checkSimulationIsParent(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := NewContext(ctx)

	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	if sim.Kind() == simulations.SimParent {
		_, err := simCtx.Services().Simulations().Reject(gid)
		if err != nil {
			return nil, err
		}
		return nil, simulations.ErrIncorrectKind
	}

	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------

// JobCheckParentSimulationWithError is used to check if a parent simulation of a certain children simulation has an error.
var JobCheckParentSimulationWithError = &actions.Job{
	Name:       "check-parent-simulation-with-error",
	Execute:    checkParentSimulationWithError,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// checkParentSimulationWithError is the main process executed by JobCheckParentSimulationWithError.
func checkParentSimulationWithError(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := NewContext(ctx)

	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	if sim.Kind() != simulations.SimChild {
		return gid, nil
	}
	parent, err := simCtx.Services().Simulations().GetParent(gid)
	if err != nil {
		return nil, err
	}

	if simerr := parent.Error(); simerr != nil {
		return nil, simulations.ErrParentSimulationWithError
	}

	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------

// JobUpdateSimulationStatusToLaunchNodes is used to set a simulation status to launch nodes.
var JobUpdateSimulationStatusToLaunchNodes = &actions.Job{
	Name:            "set-simulation-status-launch-nodes",
	Execute:         updateSimulationStatusToLaunchNodes,
	RollbackHandler: rollbackUpdateSimulationStatusToPending,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

// updateSimulationStatusToLaunchNodes is the main process executed by JobUpdateSimulationStatusToLaunchNodes.
func updateSimulationStatusToLaunchNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := NewContext(ctx)

	err := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusLaunchingNodes)
	if err != nil {
		return nil, err
	}
	return gid, nil
}

// rollbackUpdateSimulationStatusToPending is in charge of setting the status pending to the simulation
// changed by updateSimulationStatusToLaunchNodes.
func rollbackUpdateSimulationStatusToPending(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := NewContext(ctx)

	revertErr := simCtx.Services().Simulations().UpdateStatus(gid, simulations.StatusPending)
	if revertErr != nil {
		return nil, err
	}
	return gid, nil
}

//----------------------------------------------------------------------------------------------------------------------

// JobLaunchNodes is used to launch the required nodes to run a simulation.
var JobLaunchNodes = &actions.Job{
	Name:            "launch-nodes",
	PreHooks:        []actions.JobFunc{createMachineInputs, checkNodeAvailability},
	Execute:         launchNodes,
	RollbackHandler: rollbackLaunchNodes,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(map[string]interface{}{}),
}

// createMachineInputs creates the needed input for launchNodes.
func createMachineInputs(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := NewContext(ctx)

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

// checkNodeAvailability checks if the required amount of machines are available to be launched.
func checkNodeAvailability(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {
	simCtx := NewContext(ctx)

	inputMap := value.(map[string]interface{})
	createMachineInputs, ok := inputMap["createMachinesInputs"].([]cloud.CreateMachinesInput)
	if !ok {
		return nil, simulations.ErrInvalidInput
	}

	var minRequested int
	for _, c := range createMachineInputs {
		minRequested += int(c.MinCount)
	}

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

	for minRequested < simCtx.Platform().Store().Machines().Limit()-reserved {
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
		time.Sleep(simCtx.Platform().Store().Machines().PollFrequency())
	}

	return value, nil
}

// launchNodes is the main process executed by JobLaunchNodes.
func launchNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	simCtx := NewContext(ctx)

	inputMap := value.(map[string]interface{})

	gid, ok := inputMap["groupID"].(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	createMachineInputs, ok := inputMap["createMachinesInputs"].([]cloud.CreateMachinesInput)
	if !ok {
		return nil, simulations.ErrInvalidInput
	}

	sim, ok := inputMap["simulation"].(simulations.Simulation)
	if !ok {
		return nil, simulations.ErrInvalidInput
	}

	output, err := simCtx.Platform().Machines().Create(createMachineInputs)
	if len(output) != 0 && err != nil {
		if dataErr := deployment.SetJobData(tx, &deployment.CurrentJob, "machine-list", output); dataErr != nil {
			return nil, dataErr
		}
		return nil, err
	}

	err = simCtx.Services().Simulations().Update(gid, sim)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"groupID": gid,
	}, nil
}

// rollbackLaunchNodes is the process in charge of terminating the machine instances that were created in launchNodes.
func rollbackLaunchNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	if err != cloud.ErrUnknown && err != cloud.ErrInsufficientMachines && err != cloud.ErrRequestsLimitExceeded {
		return nil, nil
	}

	jobData, dataErr := deployment.GetJobData(tx, &deployment.CurrentJob, "machine-list")
	if dataErr != nil {
		return nil, err
	}

	machineList, ok := jobData.([]cloud.CreateMachinesOutput)
	if !ok {
		return nil, simulations.ErrInvalidInput
	}

	simCtx := NewContext(ctx)
	for _, m := range machineList {
		if m.ToTerminateMachinesInput() != nil {
			_ = simCtx.Platform().Machines().Terminate(*m.ToTerminateMachinesInput())
		}
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------

// JobWaitForNodes is used to wait until all required nodes are ready.
var JobWaitForNodes = &actions.Job{
	Name:       "wait-for-nodes",
	Execute:    waitForNodes,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// waitForNodes is the main process executed by JobWaitForNodes.
func waitForNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	simCtx := NewContext(ctx)

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	res := orchestrator.NewResource("", "", orchestrator.NewSelector(map[string]string{
		"cloudsim_groupid": string(gid),
	}))

	req := simCtx.Platform().Orchestrator().Nodes().WaitForCondition(res, orchestrator.ReadyCondition)

	timeout := simCtx.Platform().Store().Machines().Timeout()
	pollFreq := simCtx.Platform().Store().Machines().PollFrequency()

	err := req.Wait(timeout, pollFreq)
	if err != nil {
		return nil, err
	}

	return gid, nil
}
