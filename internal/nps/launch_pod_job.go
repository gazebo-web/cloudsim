package nps

// This file implements the launch pod job. 

import (
	"fmt"
	"strings"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

/////////////////////////////////////////////
var LaunchGazeboServerPod = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-gzserver-pod",
	PreHooks:        []actions.JobFunc{prepareGazeboCreatePodInput},
	// PostHooks:       []actions.JobFunc{},
	// RollbackHandler: rollbackPodCreation,
	InputType:       actions.GetJobDataType(&StartSimulationData{}),
	OutputType:      actions.GetJobDataType(&StartSimulationData{}),
})

func prepareGazeboCreatePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
  fmt.Printf("\n\nCreating Pods!\n\n")

	startData := store.State().(*StartSimulationData)

  var sim Simulation
  if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&sim).Error; err != nil {
    return nil, err
  }


	// What is this, and why is it needed???
	// namespace := startData.Platform().Store().Orchestrator().Namespace()
  fmt.Printf("-------------------------\n")

	// TODO: Get ports from Ignition Store
	ports := []int32{11345, 11311}

	// Set up container configuration
	privileged := true
	allowPrivilegeEscalation := true

	volumes := []orchestrator.Volume{
		{
			Name:      "xauth",
			MountPath: "/tmp/.docker.xauth",
			HostPath:  "/tmp/.docker.xauth",
		},
		{
			Name:      "localtime",
			MountPath: "/etc/localtime",
			HostPath:  "/etc/localtime",
		},
		{
			Name:      "devinput",
			MountPath: "/dev/input",
			HostPath:  "/dev/input",
		},
		{
			Name:      "x11",
			MountPath: "/tmp/.X11-unix",
			HostPath:  "/tmp/.X11-unix",
		},
	}

	envVars := map[string]string{
		"DISPLAY":          ":0",
		"TERM":             "",
		"QT_X11_NO_MITSHM": "1",
		"XAUTHORITY":       "/tmp/.docker.xauth",
		"USE_XVFB":         "1",
	}

  // \todo: Are the regular nameservers? Are they manadatory?
  nameservers := startData.Platform().Store().Orchestrator().Nameservers()

	return jobs.LaunchPodsInput{
		{
      // Name is the name of the pod that will be created.
      // \todo: Should this be unique, and where is name used?
			Name:                          sim.Name,

      // Namespace is the namespace where the pod will live in.
      // \todo: What is a namespace?
			Namespace:                     "web-cloudsim-integration",

      // Labels are the map of labels that will be applied to the pod.
      // \todo: What are the labels used for?
      Labels:                        map[string]string{"key":"value"},

      // RestartPolicy defines how the pod should react after an error.
      // \todo: What are the restart policies, and how do I choose one?
			RestartPolicy:                 orchestrator.RestartPolicyNever,

      // TerminationGracePeriodSeconds is the time duration in seconds the pod needs to terminate gracefully.
      // \todo: What does this do?
			TerminationGracePeriodSeconds: 0,

      // NodeSelector defines the node where the pod should run in.
      // \todo: What does this mean, and how do I know what value to put in???
			NodeSelector:                  orchestrator.NewSelector(map[string]string{
    "cloudsim_groupid": startData.GroupID.String() }),

      // Containers is the list of containers that should be created inside the pod.
      // \todo: What is a container? 
			Containers: []orchestrator.Container{
        {
          // Name is the container's name.
					Name:                     sim.Name,
          // Image is the image running inside the container.
					Image:                    sim.Image,
          // Args passed to the Command. Cannot be updated.
					Args:                     strings.Split(sim.Args,","),
          // Privileged defines if the container should run in privileged mode.
					Privileged:               &privileged,
          // AllowPrivilegeEscalation is used to define if the container is allowed to scale its privileges.
					AllowPrivilegeEscalation: &allowPrivilegeEscalation,
          // Ports is the list of ports that should be opened.
					Ports:                    ports,
          // Volumes is the list of volumes that should be mounted in the container.
					Volumes:                  volumes,
          // EnvVars is the list of env vars that should be passed into the container.
					EnvVars:                  envVars,
				},
			},
			Volumes:     volumes,

      // \todo: Is this required?
			Nameservers: nameservers,
		},
	}, nil
}
