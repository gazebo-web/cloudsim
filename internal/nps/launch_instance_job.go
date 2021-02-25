package nps

// This file implements the launch instance job. 

import (
  "fmt"
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
  fmt.Printf("\n\nLaunching!\n\n")
	startData := value.(*StartSimulationData)

  // \todo this doesn't return the correct zone. It looks like a subnet 
  // is returned in both parameters.
	subnet, zone := startData.Platform().Store().Machines().SubnetAndZone()
  zone = "us-east-1c"

  // \todo What is this? This line segfaults.
	/*sim, err := startData.Services().Simulations().Get(startData.GroupID)
	if err != nil {
		return nil, err
	}
  */

  initScript := "bash"
  tags := []cloud.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"Name":                       "nps-name",
				"cloudsim_groupid":           string(startData.GroupID),
				"CloudsimGroupID":            string(startData.GroupID),
				"project":                    "nps",
				"Cloudsim":                   "True",
				"cloudsim-application":       "nps",
				// "cloudsim-simulation-worker": m.NamePrefixValue,
				// "cloudsim_node_type":         nodeType,
				// clusterKey:                   "owned",
			},
		},
  }

	input := []cloud.CreateMachinesInput{
		{
      // \TODO: What is this, and how do I set the value?
      // \todo: My issue with patterns like `startData.Platform().Store().Machines().InstanceProfile()` is that it's very difficult to follow.
      //
      // Figured out that this can be traced back to env.machineStore, which
      // in turn reads the following environment variable:
      //     CLOUDSIM_MACHINES_INSTANCE_PROFILE
			InstanceProfile: startData.Platform().Store().Machines().InstanceProfile(),
      // \TODO: What is this, and how do I set the value?
      //
      // This is set from the CLOUDSIM_MACHINES_KEY_NAME environment variable
			KeyName:         startData.Platform().Store().Machines().KeyName(),

      // This appears to be the Ec2 machine type.
      // \todo: How is this configured?
      //
      // This is set from the CLOUDSIM_MACHINES_TYPE environment variable
			Type:            startData.Platform().Store().Machines().Type(),

      // \todo: This is the AMI? Not a Docker image?
      // \todo: How is this configured?
      //
      // This is set from the CLOUDSIM_MACHINES_BASE_IMAGE environment variable
			Image:           startData.Platform().Store().Machines().BaseImage(),

      // \todo: What is this?
			MinCount:        1,

      // \todo: What is this?
			MaxCount:        1,

      // \todo: How is this configured?
      //
      // This is set from the CLOUDSIM_MACHINES_FIREWALL_RULES environment
      // variable
			FirewallRules:   startData.Platform().Store().Machines().FirewallRules(),

      // \todo: What is this and how is this configured?
			SubnetID:        subnet,

      // \todo: What is this and how is this configured?
			Zone:            zone,

      // \todo: What is this and how is this configured?
			//Tags:            startData.Platform().Store().Machines().Tags(sim, "gzserver", "gzserver"),
			Tags:            tags,

      // \todo: What is this and how is this configured?
			//InitScript:      startData.Platform().Store().Machines().InitScript(),
			InitScript:      &initScript,

      // \todo: What is this and what is a good value?
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
