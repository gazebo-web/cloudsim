package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/cmdgen"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// LaunchCommsBridgePods launches the list of comms bridge and copy pods.
var LaunchCommsBridgePods = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-comms-bridge-pods",
	PreHooks:        []actions.JobFunc{setStartState, prepareCommsBridgePodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackLaunchCommsBridgePods,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func rollbackLaunchCommsBridgePods(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	for i := range robots {
		name := subtapp.GetPodNameCommsBridge(s.GroupID, subtapp.GetRobotID(i))
		ns := s.Platform().Store().Orchestrator().Namespace()

		_, _ = s.Platform().Orchestrator().Pods().Delete(resource.NewResource(name, ns, nil))
	}

	return nil, nil
}

// prepareCommsBridgePodInput prepares the input for the generic LaunchPods job to launch comms bridge pods.
func prepareCommsBridgePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	subtSim := sim.(subt.Simulation)

	track, err := s.SubTServices().Tracks().Get(subtSim.GetTrack(), subtSim.GetWorldIndex(), subtSim.GetRunIndex())
	if err != nil {
		return nil, err
	}

	var podInputs []pods.CreatePodInput

	marsupials := subtSim.GetMarsupials()

	for i, r := range subtSim.GetRobots() {
		// Create comms bridge input
		privileged := true
		allowPrivilegesEscalation := true

		initVolumes := []pods.Volume{
			{
				Name:      "logs",
				HostPath:  "/tmp",
				MountPath: "/tmp",
			},
		}

		volumes := []pods.Volume{
			{
				Name:         "logs",
				HostPath:     "/tmp/robot-logs",
				HostPathType: pods.HostPathDirectoryOrCreate,
				MountPath:    s.Platform().Store().Ignition().ROSLogsPath(),
			},
		}

		args, err := cmdgen.CommsBridge(track.World, i, r.GetName(), r.GetKind(), subt.IsRobotChildMarsupial(subtSim.GetMarsupials(), r))
		if err != nil {
			return nil, err
		}

		podInputs = append(podInputs, pods.CreatePodInput{
			Name:                          subtapp.GetPodNameCommsBridge(s.GroupID, subtapp.GetRobotID(i)),
			Namespace:                     s.Platform().Store().Orchestrator().Namespace(),
			Labels:                        subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r).Map(),
			RestartPolicy:                 pods.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
			InitContainers: []pods.Container{
				pods.NewChownContainer(initVolumes),
			},
			Containers: []pods.Container{
				{
					Name:                     subtapp.GetContainerNameCommsBridge(),
					Image:                    track.BridgeImage,
					Args:                     args,
					Privileged:               &privileged,
					AllowPrivilegeEscalation: &allowPrivilegesEscalation,
					Volumes:                  volumes,
					EnvVars: subtapp.GetEnvVarsCommsBridge(
						s.GroupID,
						r.GetName(),
						s.GazeboServerIP,
						s.Platform().Store().Ignition().Verbosity(),
					),
					EnvVarsFrom: subtapp.GetEnvVarsFromSourceCommsBridge(),
				},
			},
			Volumes:     volumes,
			Nameservers: s.Platform().Store().Orchestrator().Nameservers(),
		})

	}

	return jobs.LaunchPodsInput(podInputs), nil
}
