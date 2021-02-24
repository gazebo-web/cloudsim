package nps

// This file implements the launch instance job. 

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// LaunchInstances is a job that is used to launch machine instances.
var LaunchInstances = jobs.LaunchInstances.Extend(actions.Job{
	Name:       "launch-instances",
	PreHooks:   []actions.JobFunc{createLaunchInstancesInput},
	//PostHooks:  []actions.JobFunc{checkLaunchInstancesOutput, saveLaunchInstancesOutput, returnState},
	InputType:  actions.GetJobDataType(&StartSimulationData{}),
	OutputType: actions.GetJobDataType(&StartSimulationData{}),
})

// createLaunchInstancesInput is in charge of populating the data for the generic LaunchInstances job input.
func createLaunchInstancesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	startData := value.(*StartSimulationData)

	subnet, zone := startData.Platform().Store().Machines().SubnetAndZone()
	sim, err := startData.Services().Simulations().Get(startData.GroupID)
	if err != nil {
		return nil, err
	}

	input := []cloud.CreateMachinesInput{
		{
			InstanceProfile: startData.Platform().Store().Machines().InstanceProfile(),
			KeyName:         startData.Platform().Store().Machines().KeyName(),
			Type:            startData.Platform().Store().Machines().Type(),
			Image:           startData.Platform().Store().Machines().BaseImage(),
			MinCount:        1,
			MaxCount:        1,
			FirewallRules:   startData.Platform().Store().Machines().FirewallRules(),
			SubnetID:        subnet,
			Zone:            zone,
			Tags:            startData.Platform().Store().Machines().Tags(sim, "gzserver", "gzserver"),
			InitScript:      startData.Platform().Store().Machines().InitScript(),
			Retries:         10,
		},
	}

  // \todo: Is this needed?
	// startData.CreateMachinesInput = input
	store.SetState(startData)
	return jobs.LaunchInstancesInput(input), nil
}

// checkLaunchInstancesOutput checks that the requested instances matches the amount of created instances.
/*func checkLaunchInstancesOutput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.LaunchInstancesOutput)
	s := store.State().(*state.StartSimulation)
	var requested int64
	for _, c := range s.CreateMachinesInput {
		requested += c.MinCount
	}
	var created int64
	for _, c := range out {
		created += int64(len(c.Instances))
	}
	if requested > created {
		return nil, fmt.Errorf("not enough machines created, requested: [%d] created: [%d]", requested, created)
	}
	return out, nil
}
// saveLaunchInstancesOutput saves the list of machines created in the store.
func saveLaunchInstancesOutput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.LaunchInstancesOutput)
	s := store.State().(*state.StartSimulation)
	s.CreateMachinesOutput = out
	store.SetState(s)
	return s, nil
}*/
