package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/gazebo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

const dataGzServerCreationKey = "gz-server-creation"

// LaunchGazeboServerPod is used to launch the gazebo server pods for a simulation.
var LaunchGazeboServerPod = &actions.Job{
	Name:            "launch-gazebo-server-pod",
	PreHooks:        []actions.JobFunc{prepareGazeboServerPodConfig},
	Execute:         launchGazeboServerPod,
	RollbackHandler: rollbackLaunchGazeboServerPod,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

func prepareGazeboServerPodConfig(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
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
	// namespace := simCtx.Platform().Store().Cluster().Namespace()

	// Set up node selector
	nodeSelector := orchestrator.NewSelector(map[string]string{
		// TODO: Make keys constant
		"cloudsim_groupid":   string(gid),
		"cloudsim_node_type": "gzserver",
	})

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

	subtSim, ok := sim.(subt.Simulation)

	// Get track name
	trackName := subtSim.Track()
	subtSimService := simCtx.Services().(application.Services)
	track, err := subtSimService.Tracks().Get(trackName)
	if err != nil {
		return nil, err
	}

	// Assign track's image as container image.
	containerImage := track.Image

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

	// Generate gazebo command args
	runCommand, err := gazebo.GenerateLaunchArgs(gazebo.LaunchParams{
		Worlds:                  "",
		WorldMaxSimSeconds:      "",
		Seeds:                   nil,
		RunIndex:                nil,
		AuthorizationToken:      nil,
		MaxWebsocketConnections: 0,
		Robots:                  nil,
		Marsupials:              nil,
	})
	if err != nil {
		return nil, err
	}

	privileged := true
	allowPrivilegeEscalation := true

	ports := []int32{11345, 11311}

	volumes := []orchestrator.Volume{
		{
			Name:      "logs",
			MountPath: simCtx.Platform().Store().Ignition().LogsMountPath(),
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
		"IGN_RELAY":        simCtx.Platform().Store().Ignition().IP(), // IP Cloudsim
		"IGN_PARTITION":    string(gid),
		"IGN_VERBOSE":      simCtx.Platform().Store().Ignition().Verbosity(),
	}

	nameservers := simCtx.Platform().Store().Orchestrator().Nameservers()

	// Create the input for the operation
	input := orchestrator.CreatePodInput{
		Name:                          podName,
		Namespace:                     namespace,
		Labels:                        labels,
		RestartPolicy:                 orchestrator.RestartPolicyNever,
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
	return map[string]interface{}{
		"groupID":        gid,
		"createPodInput": input,
		"labels":         labels,
		"podName":        podName,
		"namespace":      namespace,
	}, nil
}

func launchGazeboServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Create ctx
	simCtx := context.NewContext(ctx)

	// Parse input
	inputMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	createPodInput, ok := inputMap["createPodInput"].(orchestrator.CreatePodInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	gid, ok := inputMap["groupID"].(simulations.GroupID)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	// Create pod
	_, err := simCtx.Platform().Orchestrator().Pods().Create(createPodInput)
	if dataErr := deployment.SetJobData(tx, nil, dataGzServerCreationKey, value); dataErr != nil {
		return nil, dataErr
	}
	if err != nil {
		return nil, err
	}

	return gid, nil
}

func rollbackLaunchGazeboServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {

	// Create ctx
	simCtx := context.NewContext(ctx)

	// Get pod name
	jobCreationData, dataErr := deployment.GetJobData(tx, &deployment.CurrentJob, dataGzServerCreationKey)
	if dataErr != nil {
		return nil, dataErr
	}

	// Parse input
	inputMap, ok := jobCreationData.(map[string]interface{})
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	name, ok := inputMap["podName"].(string)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	namespace, ok := inputMap["namespace"].(string)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	labels, ok := inputMap["labels"].(map[string]string)
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
