package nps

// This file implements the launch instance job.

import (
	"encoding/base64"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"strings"
)

// LaunchInstances is a job that is used to launch machine instances.
var LaunchInstances = jobs.LaunchInstances.Extend(actions.Job{
	Name:       "launch-instances",
	PreHooks:   []actions.JobFunc{createLaunchInstancesInput},
	PostHooks:  []actions.JobFunc{checkLaunchInstancesOutput, saveLaunchInstancesOutput, returnState},
	InputType:  actions.GetJobDataType(&StartSimulationData{}),
	OutputType: actions.GetJobDataType(&StartSimulationData{}),
})

// createLaunchInstancesInput is in charge of populating the data for the generic LaunchInstances job input.
func createLaunchInstancesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Get the start simulation data for this job.
	startData := value.(*StartSimulationData)

	// Update the database entry with the latest status
	// \todo Help needed: I think this is not the recommended method to update
	// the database.
	var simEntry Simulation
	if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&simEntry).Error; err != nil {
		return nil, err
	}
	simEntry.Status = "launching"
	tx.Save(&simEntry)

	// Get the subtnet and zone to use. This will return bad values unless
	// you set the correct environment variables.
	// \todo Improvement: Return an error if the environment variables have not been set.
	subnet, zone := startData.Platform().Store().Machines().SubnetAndZone()

	// This is a magic command that lets the EC2 machine join the Kubernetes
	// cluster.
	// \todo Improvement: Make this easier to find and customize.
	command := `
  #!/bin/bash
  set -x
  set -o xtrace
  cat > /etc/systemd/system/kubelet.service.d/20-labels-taints.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS=--node-labels=cloudsim_groupid=` + startData.GroupID.String() + `"
EOF

  /etc/eks/bootstrap.sh %s %s
`

	// \todo Help needed: Copied from Subt. Is this needed?
	arguments := []string{
		// Allow the node to contain unlimited pods
		"--use-max-pods false",
	}

	// The cluster name is read from the CLOUDSIM_MACHINES_CLUSTER_NAME
	// environment variable.
	clusterName := startData.Platform().Store().Machines().ClusterName()
	initScript := fmt.Sprintf(command, clusterName,
		strings.Join(arguments, " "))
	initScript = base64.StdEncoding.EncodeToString([]byte(initScript))

	clusterKey := "kubernetes.io/cluster/" + clusterName

	// These are the tags to apply the EC2 machines.
	tags := []cloud.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"Name":                 simEntry.Name,
				"cloudsim_groupid":     string(startData.GroupID),
				"project":              "nps",
				"Cloudsim":             "true",
				"cloudsim-application": "nps",

				// Note: `clusterKey` is extremely important. Without it, the EC2 node
				// will not join the cluster.
				clusterKey: "owned",
			},
		},
	}

	input := []cloud.CreateMachinesInput{
		{
			// This value can be traced back to env.machineStore, which
			// in turn reads the following environment variable:
			//     CLOUDSIM_MACHINES_INSTANCE_PROFILE
			InstanceProfile: startData.Platform().Store().Machines().InstanceProfile(),
			// This is set from the CLOUDSIM_MACHINES_KEY_NAME environment variable
			KeyName: startData.Platform().Store().Machines().KeyName(),

			// This is the EC2 machine type and is set from the
			// CLOUDSIM_MACHINES_TYPE environment variable
			Type: startData.Platform().Store().Machines().Type(),

			// This is set from the CLOUDSIM_MACHINES_BASE_IMAGE environment variable
			// This is the AMI that is loaded onto an EC2 instances.
			Image: startData.Platform().Store().Machines().BaseImage(),

			// \todo Help needed: What is this?
			MinCount: 1,

			// \todo Help needed: What is this?
			MaxCount: 1,

			// This is set from the CLOUDSIM_MACHINES_FIREWALL_RULES environment
			// variable
			FirewallRules: startData.Platform().Store().Machines().FirewallRules(),

			// This is set from the CLOUDSIM_MACHINES_SUBNETS
			SubnetID: subnet,

			// This is set from the CLOUDSIM_MACHINES_ZONES
			Zone: zone,

			// Tags to apply to the EC2 machine.
			Tags: tags,

			// This init script is the command to run when the base image is loaded
			// onto the EC2 machine.
			InitScript: &initScript,

			// \todo Help needed: What is this and what is a good value?
			Retries: 1,
		},
	}

	// This will allow us to select the node launched by this job in future
	// jobs. See --node-labels=cloudsim_groupid=GROUP_ID in the InitScript above.
	startData.NodeSelector = orchestrator.NewSelector(map[string]string{
		"cloudsim_groupid": startData.GroupID.String(),
	})

	// \todo Help needed: What is this, and why should I store it?
	startData.CreateMachinesInput = input
	store.SetState(startData)

	return jobs.LaunchInstancesInput(input), nil
}

// checkLaunchInstancesOutput checks that the requested instances matches the amount of created instances.
//
// \todo Improvement: I directly copied this from SubT. Can this be captured as a general purpose function?
func checkLaunchInstancesOutput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.LaunchInstancesOutput)
	s := store.State().(*StartSimulationData)
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
// \todo Improvement: I directly copied this from SubT. Can this be captured as a general purpose function?
func saveLaunchInstancesOutput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.LaunchInstancesOutput)
	s := store.State().(*StartSimulationData)
	s.CreateMachinesOutput = out
	store.SetState(s)
	return s, nil
}
