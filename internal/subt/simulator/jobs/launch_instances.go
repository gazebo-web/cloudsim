package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// LaunchInstances is a job that is used to launch machine instances for simulations.
var LaunchInstances = jobs.LaunchInstances.Extend(actions.Job{
	Name:            "launch-instances",
	PreHooks:        []actions.JobFunc{setStartState, createLaunchInstancesInput, checkMachinesLimit},
	PostHooks:       []actions.JobFunc{checkLaunchInstancesOutput, saveLaunchInstancesOutput, returnState},
	RollbackHandler: removeLaunchedInstances,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// checkMachinesLimit is used to check if there are enough machines before proceeding with requesting them.
func checkMachinesLimit(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	// if no limit is set, skip the pre-hook execution.
	if s.Platform().Store().Machines().Limit() < 0 {
		return value, nil
	}

	input := value.(jobs.LaunchInstancesInput)

	// Count the amount of requested machines.
	var requested int
	for _, in := range input {
		requested += int(in.MaxCount)
	}

	// Get the number of provisioned machines from cloud provider.
	reserved := s.Platform().Machines().Count(machines.CountMachinesInput{
		Filters: map[string][]string{
			"instance-state-name": {
				"pending",
				"running",
			},
		},
	})

	// Check that the number of required machines is greater than the current available amount of machines.
	if requested > s.Platform().Store().Machines().Limit()-reserved {
		return nil, machines.ErrInsufficientMachines
	}

	return value, nil
}

func removeLaunchedInstances(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
	tags := subtapp.GetTagsInstanceBase(s.GroupID)

	filters := make(map[string][]string)

	for _, tag := range tags {
		for k, v := range tag.Map {
			filters[fmt.Sprintf("tag:%s", k)] = []string{v}
		}
	}

	_ = s.Platform().Machines().Terminate(machines.TerminateMachinesInput{
		Filters: filters,
	})

	return nil, nil
}

// createLaunchInstancesInput is in charge of populating the data for the generic LaunchInstances job input.
func createLaunchInstancesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
	subnet, zone := s.Platform().Store().Machines().SubnetAndZone()

	prefix := s.Platform().Store().Machines().NamePrefix()
	clusterName := s.Platform().Store().Machines().ClusterName()

	input := []machines.CreateMachinesInput{
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
			Tags:            subtapp.GetTagsInstanceSpecific(prefix, s.GroupID, "gzserver", clusterName, "gzserver"),
			Retries:         10,
			Labels:          subtapp.GetNodeLabelsGazeboServer(s.GroupID).Map(),
			ClusterID:       clusterName,
		},
	}

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	for _, r := range robots {
		tags := subtapp.GetTagsInstanceSpecific(prefix, s.GroupID, fmt.Sprintf("fc-%s", r.GetName()), clusterName, "field-computer")

		input = append(input, machines.CreateMachinesInput{
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
			Retries:         10,
			Labels:          subtapp.GetNodeLabelsFieldComputer(s.GroupID, r).Map(),
			ClusterID:       clusterName,
		})
	}

	s.CreateMachinesInput = input

	store.SetState(s)

	return jobs.LaunchInstancesInput(input), nil
}

// checkLaunchInstancesOutput checks that the requested instances matches the amount of created instances.
func checkLaunchInstancesOutput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
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
}
