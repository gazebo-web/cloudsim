package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/cmdgen"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"path"
)

// LaunchCommsBridgePods launches the list of comms bridge and copy pods.
var LaunchCommsBridgePods = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-comms-bridge-pods",
	PreHooks:        []actions.JobFunc{setStartState, prepareCommsBridgePodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackPodCreation,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

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

	var pods []orchestrator.CreatePodInput

	for i, r := range subtSim.GetRobots() {
		hostPath := "/tmp"
		logDirectory := "robot-logs"
		logMountPath := path.Join(hostPath, logDirectory)

		// Create comms bridge input
		privileged := true
		allowPrivilegesEscalation := true

		volumes := []orchestrator.Volume{
			{
				Name:         "logs",
				HostPath:     logMountPath,
				HostPathType: orchestrator.HostPathDirectoryOrCreate,
				MountPath:    s.Platform().Store().Ignition().ROSLogsPath(),
			},
		}

		args, err := cmdgen.CommsBridge(track.World, i, r.GetName(), r.GetKind(), subt.IsRobotChildMarsupial(subtSim.GetMarsupials(), r))
		if err != nil {
			return nil, err
		}

		pods = append(pods, orchestrator.CreatePodInput{
			Name:                          subtapp.GetPodNameCommsBridge(s.GroupID, subtapp.GetRobotID(i)),
			Namespace:                     s.Platform().Store().Orchestrator().Namespace(),
			Labels:                        subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r).Map(),
			RestartPolicy:                 orchestrator.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
			InitContainers: []orchestrator.Container{
				orchestrator.NewChownContainer(volumes),
			},
			Containers: []orchestrator.Container{
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

	return jobs.LaunchPodsInput(pods), nil
}
