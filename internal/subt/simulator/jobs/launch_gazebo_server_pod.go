package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/gazebo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// LaunchGazeboServerPod is used to launch the gazebo server pods for a simulation.
var LaunchGazeboServerPod = LaunchPod.Extend(actions.Job{
	Name:            "launch-gazebo-server-pod",
	PreHooks:        []actions.JobFunc{prepareGazeboServerPodConfig},
	PostHooks:       []actions.JobFunc{launchGazeboServerPodPostHook},
	RollbackHandler: rollbackLaunchGazeboServerPod,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
})

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

	// Get simulation
	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	// Parse to subt simulation
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

	// Get terminate grace period
	terminationGracePeriod := simCtx.Platform().Store().Orchestrator().TerminationGracePeriod()

	// Generate gazebo command args
	// TODO: Fill LaunchParams with requested information
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

	// Set up container configuration
	privileged := true
	allowPrivilegeEscalation := true

	// TODO: Get ports from Ignition Store
	ports := []int32{11345, 11311}

	volumes := []orchestrator.Volume{
		{
			Name:      "logs",
			MountPath: simCtx.Platform().Store().Ignition().GazeboServerLogsPath(),
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

	data := simCtx.Store().Get().(*StartSimulationData)

	nodeSelector := orchestrator.NewSelector(data.GazeboNodeSelector)

	// Create the input for the operation
	input := orchestrator.CreatePodInput{
		Name:                          podName,
		Namespace:                     namespace,
		Labels:                        data.GazeboLabels,
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
	return input, nil
}

func launchGazeboServerPodPostHook(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	data := ctx.Store().Get().(*StartSimulationData)

	output := value.(LaunchPodOutput)

	// Save resource
	data.GazeboPodResource = output.Resource

	dataErr := ctx.Store().Set(data)
	if dataErr != nil {
		return nil, dataErr
	}

	// Check error
	if output.Error != nil {
		return nil, output.Error
	}

	return data.GroupID, nil
}

func rollbackLaunchGazeboServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	simCtx := context.NewContext(ctx)

	data := simCtx.Store().Get().(*StartSimulationData)

	_, delErr := simCtx.Platform().Orchestrator().Pods().Delete(data.GazeboPodResource)
	if delErr != nil {
		return nil, delErr
	}

	return nil, nil
}
