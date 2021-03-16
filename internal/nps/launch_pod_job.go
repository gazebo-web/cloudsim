package nps

// This file implements the launch pod job.

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"strings"
)

/////////////////////////////////////////////
var LaunchPod = jobs.LaunchPods.Extend(actions.Job{
	Name:      "launch-pod",
	PreHooks:  []actions.JobFunc{setStartState, prepareCreatePodInput},
	PostHooks: []actions.JobFunc{returnState},
	// RollbackHandler: rollbackPodCreation,
	InputType:  actions.GetJobDataType(&StartSimulationData{}),
	OutputType: actions.GetJobDataType(&StartSimulationData{}),
})

func prepareCreatePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {

	startData := store.State().(*StartSimulationData)

	// Update the database entry with the latest status
	// \todo Help needed: I think this is not the recommended method to update
	// the database.
	var sim Simulation
	if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&sim).Error; err != nil {
		return nil, err
	}
	sim.Status = "Creating docker image (pod)."
	tx.Save(&sim)

	// Namespace is the orchestrator namespace where simulations should be
	// launched.
	// \todo MAJOR ERROR: I would assume that this would return the value in
	// CLOUDSIM_MACHINES_ORCHESTRATOR_NAMESPACE. It is empty.
	namespace := startData.Platform().Store().Orchestrator().Namespace()
	if namespace == "default" || namespace == "" {
		startData.logger.Error("In prepareCreatePodInput, CLOUDSIM_ORCHESTRATOR_NAMESPACE has not been set")
		return nil, errors.New("CLOUDSIM_ORCHESTRATOR_NAMESPACE has not been set")
	}
	// namespace := "web-cloudsim-integration"

	// \todo Improvment: Get ports dynamically.
	ports := []orchestrator.ContainerPort{
		orchestrator.ContainerPort{
			ContainerPort: 8080,
			HostPort:      8080,
		},
	}

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
		"DISPLAY":          "",
		"QT_X11_NO_MITSHM": "1",
		"XAUTHORITY":       "/tmp/.docker.xauth",
		"USE_XVFB":         "1",
	}

	// \todo Help needed: Are the regular nameservers? Are they manadatory?
	nameservers := startData.Platform().Store().Orchestrator().Nameservers()
	labels := map[string]string{
		"cloudsim":         "true",
		"nps":              "true",
		"cloudsim_groupid": startData.GroupID.String(),
	}

	startData.PodSelector = orchestrator.NewSelector(labels)

	var args []string
	if sim.Args != "" {
		args = strings.Split(sim.Args, ",")
	}
	startData.logger.Info("Launching pod. Image[", sim.Image, "] Args[", args, "]")
	return jobs.LaunchPodsInput{
		{
			// Name is the name of the pod that will be created.
			// \todo: Should this be unique, and where is name used?
			Name: startData.GroupID.String(),

			// Namespace is the namespace where the pod will live in.
			Namespace: namespace,

			// Labels are the map of labels that will be applied to the pod.
			// These labels are very important in order to reference
			// kubernetes resources in other places, just as jobs.
			Labels: labels,

			// RestartPolicy defines how the pod should react after an error.
			// \todo Help needed: What are the restart policies, and how do I
			// choose one?
			RestartPolicy: orchestrator.RestartPolicyNever,

			// TerminationGracePeriodSeconds is the time duration in seconds the pod needs to terminate gracefully.
			// \todo Help needed: What does this do?
			TerminationGracePeriodSeconds: 0,

			// NodeSelector defines the node where the pod should run in. This is
			// very important and must be set correctly to match a label on a Node
			// otherwise the pod will not run and remain in a `pending` state.
			//
			// A Node's labels are set when launching an instance via
			// `KUBELET_EXTRA_ARGS=--node-labels=KEY=VALUE` in a
			//   `/etc/systemd/system/kubelet.service.d/20-labels-taints.conf` file.
			NodeSelector: startData.NodeSelector,

			// Containers is the list of containers that should be created inside the pod.
			Containers: []orchestrator.Container{
				{
					// Name is the container's name.
					Name: startData.GroupID.String(),
					// Image is the image running inside the container.
					Image: sim.Image,
					// Args passed to the Command. Cannot be updated.
					Args: args,
					// Privileged defines if the container should run in privileged mode.
					Privileged: &privileged,
					// AllowPrivilegeEscalation is used to define if the container is allowed to scale its privileges.
					AllowPrivilegeEscalation: &allowPrivilegeEscalation,
					// Ports is the list of ports that should be opened.
					Ports: ports,
					// Volumes is the list of volumes that should be mounted in the container.
					Volumes: volumes,
					// EnvVars is the list of env vars that should be passed into the container.
					EnvVars: envVars,
				},
			},
			Volumes: volumes,

			Nameservers: nameservers,

			HostNetwork: true,
		},
	}, nil
}
