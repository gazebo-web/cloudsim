package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// LaunchInstances is a job that is used to launch machine instances for simulations.
var LaunchInstances = jobs.LaunchInstances.Extend(actions.Job{
	Name:       "launch-instances",
	PreHooks:   []actions.JobFunc{setStartState, createLaunchInstancesInput},
	PostHooks:  []actions.JobFunc{checkLaunchInstancesOutput, saveLaunchInstancesOutput, returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createLaunchInstancesInput is in charge of populating the data for the generic LaunchInstances job input.
func createLaunchInstancesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)
	subnet, zone := s.Platform().Store().Machines().SubnetAndZone()

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	input := []cloud.CreateMachinesInput{
		{
			InstanceProfile: s.Platform().Store().Machines().InstanceProfile(),
			KeyName:         s.Platform().Store().Machines().KeyName(),
			Type:            s.Platform().Store().Machines().Type(),
			Image:           s.Platform().Store().Machines().BaseImage(),
			MinCount:        1,
			MaxCount:        1,
			FirewallRules:   s.Platform().Store().Machines().FirewallRules(),
			SubnetID:        subnet,
			Zone:            zone,
			Tags:            s.Platform().Store().Machines().Tags(sim, "gzserver", "gzserver"),
			InitScript:      s.Platform().Store().Machines().InitScript(),
			Retries:         10,
		},
	}

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	for _, r := range robots {
		tags := s.
			Platform().
			Store().
			Machines().
			Tags(sim, "field-computer", fmt.Sprintf("fc-%s", r.Name()))

		input = append(input, cloud.CreateMachinesInput{
			InstanceProfile: s.Platform().Store().Machines().InstanceProfile(),
			KeyName:         s.Platform().Store().Machines().KeyName(),
			Type:            s.Platform().Store().Machines().Type(),
			Image:           s.Platform().Store().Machines().BaseImage(),
			MinCount:        1,
			MaxCount:        1,
			FirewallRules:   s.Platform().Store().Machines().FirewallRules(),
			SubnetID:        subnet,
			Zone:            zone,
			Tags:            tags,
			InitScript:      s.Platform().Store().Machines().InitScript(),
			Retries:         10,
		})
	}

	s.CreateMachinesInputs = input

	store.SetState(s)

	return jobs.LaunchInstancesInput(input), nil
}

// checkLaunchInstancesOutput checks that the requested instances matches the amount of created instances.
func checkLaunchInstancesOutput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.LaunchInstancesOutput)

	s := store.State().(*state.StartSimulation)

	var requested int64
	for _, c := range s.CreateMachinesInputs {
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

	s.CreateMachinesOutputs = out

	store.SetState(s)

	return s, nil
}
