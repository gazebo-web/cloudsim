package jobs

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"path"
)

// LaunchCommsBridge launches the list of comms bridge and copy pods.
var LaunchCommsBridge = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-comms-bridge-pods",
	PreHooks:        []actions.JobFunc{setStartState, prepareCommsBridgePodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackPodCreation,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func prepareCommsBridgePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {

	s := store.State().(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	subtSim := sim.(subt.Simulation)

	track, err := s.SubTServices().Tracks().Get(subtSim.GetTrack())
	if err != nil {
		return nil, err
	}

	var pods []orchestrator.CreatePodInput

	for i, r := range subtSim.GetRobots() {
		childMarsupial := "false"
		if subt.IsRobotChildMarsupial(subtSim.GetMarsupials(), r) {
			childMarsupial = "true"
		}

		hostPath := "/tmp"
		logDirectory := "robot-logs"
		logMountPath := path.Join(hostPath, logDirectory)

		// Create comms bridge input
		pods = append(pods, prepareCreatePodInput(configCreatePodInput{
			name:                   subtapp.GetPodNameCommsBridge(s.GroupID, subtapp.GetRobotID(i+1)),
			namespace:              s.Platform().Store().Orchestrator().Namespace(),
			labels:                 subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r).Map(),
			restartPolicy:          orchestrator.RestartPolicyNever,
			terminationGracePeriod: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			nodeSelector:           subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
			containerName:          subtapp.GetContainerNameCommsBridge(),
			image:                  track.BridgeImage,
			args: []string{
				track.World,
				fmt.Sprintf("robotName%d:=%s", i, r.Name()),
				fmt.Sprintf("robotConfig%d:=%s", i, r.Kind()),
				"headless:=true",
				fmt.Sprintf("marsupial:=%s", childMarsupial),
			},
			privileged:                true,
			allowPrivilegesEscalation: true,
			volumes: []orchestrator.Volume{
				{
					Name:         "logs",
					HostPath:     logMountPath,
					HostPathType: orchestrator.HostPathDirectoryOrCreate,
					MountPath:    s.Platform().Store().Ignition().ROSLogsPath(),
				},
			},
			envVars: subtapp.GetEnvVarsCommsBridge(
				s.GroupID,
				r.Name(),
				s.GazeboServerIP,
				s.Platform().Store().Ignition().Verbosity(),
			),
			nameservers: s.Platform().Store().Orchestrator().Nameservers(),
		}))

		if s.Platform().Store().Ignition().LogsCopyEnabled() {
			secretsName := s.Platform().Store().Ignition().SecretsName()
			ns := s.Platform().Store().Orchestrator().Namespace()

			secret, err := s.Platform().Secrets().Get(context.TODO(), secretsName, ns)
			if err != nil {
				return nil, err
			}

			accessKey := string(secret.Data[s.Platform().Store().Ignition().AccessKeyLabel()])
			secretAccessKey := string(secret.Data[s.Platform().Store().Ignition().SecretAccessKeyLabel()])

			pods = append(pods, prepareCreatePodInput(configCreatePodInput{
				name:                   subtapp.GetPodNameCommsBridgeCopy(s.GroupID, subtapp.GetRobotID(i+1)),
				namespace:              ns,
				labels:                 subtapp.GetPodLabelsCommsBridgeCopy(s.GroupID, s.ParentGroupID, r).Map(),
				restartPolicy:          orchestrator.RestartPolicyNever,
				terminationGracePeriod: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
				nodeSelector:           subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
				containerName:          subtapp.GetContainerNameCommsBridgeCopy(),
				image:                  "infrastructureascode/aws-cli:latest",
				command:                []string{"tail", "-f", "/dev/null"},
				volumes: []orchestrator.Volume{
					{
						Name:         "logs",
						HostPath:     logMountPath,
						MountPath:    s.Platform().Store().Ignition().ROSLogsPath(),
						HostPathType: orchestrator.HostPathDirectoryOrCreate,
					},
				},
				envVars: subtapp.GetEnvVarsCommsBridgeCopy(
					s.Platform().Store().Ignition().Region(),
					accessKey,
					secretAccessKey,
				),
				nameservers: s.Platform().Store().Orchestrator().Nameservers(),
			}))
		}
	}

	return jobs.LaunchPodsInput(pods), nil
}
