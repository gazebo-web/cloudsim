package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/cmdgen"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"time"
)

// LaunchGazeboServerPod launches a gazebo server pod.
var LaunchGazeboServerPod = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-gzserver-pod",
	PreHooks:        []actions.JobFunc{setStartState, prepareGazeboCreatePodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackLaunchGazeboServerPod,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func rollbackLaunchGazeboServerPod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := subtapp.GetPodNameGazeboServer(s.GroupID)
	ns := s.Platform().Store().Orchestrator().Namespace()

	_, _ = s.Platform().Orchestrator().Pods().Delete(resource.NewResource(name, ns, nil))

	return nil, nil
}

func prepareGazeboCreatePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
	// Set up namespace
	namespace := s.Platform().Store().Orchestrator().Namespace()

	// Get simulation
	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	// Parse to subt simulation
	subtSim := sim.(simulations.Simulation)

	// Get track
	track, err := s.SubTServices().Tracks().Get(subtSim.GetTrack(), subtSim.GetWorldIndex(), 0)

	if err != nil {
		return nil, err
	}
	// Generate gazebo command args
	runCommand := cmdgen.Gazebo(cmdgen.GazeboConfig{
		World:              track.World,
		WorldMaxSimSeconds: time.Duration(track.MaxSimSeconds) * time.Second,
		Seed:               track.Seed,
		AuthorizationToken: subtSim.GetToken(),
		// TODO: Get max connections from store.
		MaxWebsocketConnections: 500,
		Robots:                  subtSim.GetRobots(),
		Marsupials:              subtSim.GetMarsupials(),
		RosEnabled:              true,
	})

	// Set up container configuration
	privileged := true
	allowPrivilegeEscalation := true

	// TODO: Get ports from Ignition Store
	ports := []int32{11345, 11311}

	volumes := []pods.Volume{
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

	nameservers := s.Platform().Store().Orchestrator().Nameservers()

	return jobs.LaunchPodsInput{
		{
			Name:                          subtapp.GetPodNameGazeboServer(s.GroupID),
			Namespace:                     namespace,
			Labels:                        subtapp.GetPodLabelsGazeboServer(s.GroupID, s.ParentGroupID).Map(),
			RestartPolicy:                 pods.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsGazeboServer(s.GroupID),
			Containers: []pods.Container{
				{
					Name:                     subtapp.GetContainerNameGazeboServer(),
					Image:                    track.Image,
					Args:                     runCommand,
					Privileged:               &privileged,
					AllowPrivilegeEscalation: &allowPrivilegeEscalation,
					Ports:                    ports,
					Volumes:                  volumes,
					EnvVarsFrom:              subtapp.GetEnvVarsFromSourceGazeboServer(),
					EnvVars: subtapp.GetEnvVarsGazeboServer(
						s.GroupID,
						s.Platform().Store().Ignition().IP(),
						s.Platform().Store().Ignition().Verbosity(),
					),
				},
			},
			Volumes:     volumes,
			Nameservers: nameservers,
		},
	}, nil
}
