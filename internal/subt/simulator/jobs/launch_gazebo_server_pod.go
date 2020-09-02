package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// LaunchGazeboServerPod is used to launch the gazebo server pods for a simulation.
var LaunchGazeboServerPod = &actions.Job{
	Name:            "launch-gazebo-server-pod",
	Execute:         launchGazeboServerPod,
	RollbackHandler: rollbackLaunchGazeboServerPod,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

func launchGazeboServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Parse group id
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	// Create ctx
	simCtx := context.NewContext(ctx)

	// Set up pod name
	podName := "prefix-groupid-gzserver"

	// Set up namespace
	namespace := "default"
	// namespace := simCtx.Platform().Store().Platform().Namespace()

	// Set up node selector
	nodeSelector := orchestrator.NewSelector(map[string]string{})

	// Set up pod labels
	labels := map[string]string{
		"cloudsim":          "true",
		"pod-group":         podName,
		"cloudsim-group-id": string(gid),
		"gzserver":          "true",
	}

	// Get simulation
	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	// Get track name
	trackName := sim.Track()
	track, err := simCtx.Services().Tracks().GetByName(trackName)
	if err != nil {
		return nil, err
	}

	// Assign track's image as container image.
	containerImage := track.Image()

	// If simulation is child, add another label with the parent's group id.
	if sim.Kind() == simulations.SimChild {
		parent, err := simCtx.Services().Simulations().GetParent(gid)
		if err != nil {
			return nil, err
		}
		labels["parent-group-id"] = string(parent.GroupID())
	}

	// Get terminate grace period
	terminationGracePeriod := simCtx.Platform().Store().Orchestrator().TerminationGracePeriod()

	// Run gz command
	runCommand := []string{""}

	privileged := true
	allowPrivilegeEscalation := true

	ports := []int32{11345, 11311}

	volumes := []orchestrator.Volume{
		{
			Name:      "logs",
			MountPath: simCtx.Platform().Store().Gazebo().LogsMountPath(),
			HostPath:  "/tmp",
		},
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
		"QT_X11_NO_MITSHM": "1",
		"XAUTHORITY":       "/tmp/.docker.xauth",
		"USE_XVFB":         "1",
		"IGN_RELAY":        simCtx.Platform().Store().Ignition().RelayIP(),
		"IGN_PARTITION":    string(gid),
		"IGN_VERBOSE":      simCtx.Platform().Store().Ignition().Verbosity(),
	}

	nameservers := []string{
		"8.8.8.8",
		"1.1.1.1",
	}

	// Create the input for the operation
	input := orchestrator.CreatePodInput{
		Name:                          podName,
		Namespace:                     namespace,
		Labels:                        labels,
		RestartPolicy:                 kubernetes.RestartPolicyNever,
		TerminationGracePeriodSeconds: terminationGracePeriod,
		NodeSelector:                  nodeSelector,
		Containers: []orchestrator.Container{
			{
				Name:                     "gzserver-container",
				Image:                    containerImage,
				Args:                     runCommand,
				Privileged:               &privileged,
				AllowPrivilegeEscalation: &allowPrivilegeEscalation,
				Ports:   ports,
				Volumes: volumes,
				EnvVars: envVars,
			},
		},
		Volumes:     volumes,
		Nameservers: nameservers,
	}

	// Create pod
	_, err = simCtx.Platform().Orchestrator().Pods().Create(input)
	if err != nil {
		if dataErr := deployment.SetJobData(tx, nil, "gz-server-pod-labels", labels); dataErr != nil {
			return nil, dataErr
		}
		if dataErr := deployment.SetJobData(tx, nil, "gz-server-pod-name", podName); dataErr != nil {
			return nil, dataErr
		}
		if dataErr := deployment.SetJobData(tx, nil, "gz-server-pod-namespace", namespace); dataErr != nil {
			return nil, dataErr
		}
		return nil, err
	}

	return gid, nil
}

func rollbackLaunchGazeboServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	// Create ctx
	simCtx := context.NewContext(ctx)

	// Get pod name
	jobDataPodName, dataErr := deployment.GetJobData(tx, &deployment.CurrentJob, "gz-server-pod-name")
	if dataErr != nil {
		return nil, dataErr
	}

	name, ok := jobDataPodName.(string)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Get namespace
	jobDataNamespace, dataErr := deployment.GetJobData(tx, &deployment.CurrentJob, "gz-server-pod-namespace")
	if dataErr != nil {
		return nil, dataErr
	}

	namespace, ok := jobDataNamespace.(string)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Get labels
	jobDataLabels, dataErr := deployment.GetJobData(tx, &deployment.CurrentJob, "gz-server-pod-labels")
	if dataErr != nil {
		return nil, dataErr
	}

	labels, ok := jobDataLabels.(map[string]string)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	res := orchestrator.NewResource(name, namespace, orchestrator.NewSelector(labels))

	delErr := simCtx.Platform().Orchestrator().Pods().Delete(res)
	if delErr != nil {
		return nil, delErr
	}

	return nil, nil
}
