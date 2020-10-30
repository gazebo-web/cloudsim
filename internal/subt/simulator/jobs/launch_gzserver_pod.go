package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/gazebo"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"time"
)

// LaunchGazeboServerPod launches a gazebo server pod.
var LaunchGazeboServerPod = jobs.LaunchPod.Extend(actions.Job{
	Name:            "launch-gzserver-pod",
	PreHooks:        []actions.JobFunc{setStartState, prepareCreatePodInput},
	PostHooks:       []actions.JobFunc{returnState},
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func prepareCreatePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	// TODO: How do we get the pod name?
	podName := "prefix-groupid-gzserver"

	// Set up namespace
	namespace := s.Platform().Store().Orchestrator().Namespace()

	// Get simulation
	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	// Parse to subt simulation
	subtSim := sim.(simulations.Simulation)

	// Get track name
	trackName := subtSim.Track()
	app := s.Services().(subtapp.Services)
	track, err := app.Tracks().Get(trackName)
	if err != nil {
		return nil, err
	}

	// Assign track's image as container image.
	containerImage := track.Image

	// Get terminate grace period
	terminationGracePeriod := s.Platform().Store().Orchestrator().TerminationGracePeriod()


	// Generate gazebo command args
	runCommand := gazebo.Generate(gazebo.LaunchConfig{
		Worlds:                  []string{track.World},
		WorldMaxSimSeconds:      time.Duration(track.MaxSimSeconds),
		Seeds:                   track.Seed,
		AuthorizationToken:      subtSim.Token(),
		MaxWebsocketConnections: 500,
		Robots:                  subtSim.Robots(),
		Marsupials:              subtSim.Marsupials(),
	})

	// Set up container configuration
	privileged := true
	allowPrivilegeEscalation := true

	// TODO: Get ports from Ignition Store
	ports := []int32{11345, 11311}

	volumes := []orchestrator.Volume{
		{
			Name:      "logs",
			MountPath: s.Platform().Store().Ignition().GazeboServerLogsPath(),
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
		"IGN_RELAY":        s.Platform().Store().Ignition().IP(), // IP Cloudsim
		"IGN_PARTITION":    string(s.GroupID),
		"IGN_VERBOSE":      s.Platform().Store().Ignition().Verbosity(),
	}

	nameservers := s.Platform().Store().Orchestrator().Nameservers()

	nodeSelector := orchestrator.NewSelector(s.GazeboNodeLabels)

	// Create the input for the operation
	input := orchestrator.CreatePodInput{
		Name:                          podName,
		Namespace:                     namespace,
		Labels:                        s.GazeboServerPodLabels,
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
				Ports:                    ports,
				Volumes:                  volumes,
				EnvVars:                  envVars,
			},
		},
		Volumes:     volumes,
		Nameservers: nameservers,
	}

	return jobs.LaunchPodInput(input), nil
}
