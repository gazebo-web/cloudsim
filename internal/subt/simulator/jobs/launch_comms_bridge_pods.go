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

		pods = append(pods, orchestrator.CreatePodInput{
			Name:                          subtapp.GetPodNameCommsBridge(s.GroupID, subtapp.GetRobotID(i+1)),
			Namespace:                     s.Platform().Store().Orchestrator().Namespace(),
			Labels:                        subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r).Map(),
			RestartPolicy:                 orchestrator.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
			Containers: []orchestrator.Container{
				{
					Name:  subtapp.GetContainerNameCommsBridge(),
					Image: track.BridgeImage,
					Args: []string{
						track.World,
						fmt.Sprintf("robotName%d:=%s", i, r.Name()),
						fmt.Sprintf("robotConfig%d:=%s", i, r.Kind()),
						"headless:=true",
						fmt.Sprintf("marsupial:=%s", childMarsupial),
					},
					Privileged:               &privileged,
					AllowPrivilegeEscalation: &allowPrivilegesEscalation,
					Volumes:                  volumes,
					EnvVars: subtapp.GetEnvVarsCommsBridge(
						s.GroupID,
						r.Name(),
						s.GazeboServerIP,
						s.Platform().Store().Ignition().Verbosity(),
					),
				},
			},
			Volumes:     volumes,
			Nameservers: s.Platform().Store().Orchestrator().Nameservers(),
		})

		if s.Platform().Store().Ignition().LogsCopyEnabled() {
			secretsName := s.Platform().Store().Ignition().SecretsName()
			ns := s.Platform().Store().Orchestrator().Namespace()

			secret, err := s.Platform().Secrets().Get(context.TODO(), secretsName, ns)
			if err != nil {
				return nil, err
			}

			accessKey := string(secret.Data[s.Platform().Store().Ignition().AccessKeyLabel()])
			secretAccessKey := string(secret.Data[s.Platform().Store().Ignition().SecretAccessKeyLabel()])

			pods = append(pods, orchestrator.CreatePodInput{
				Name:                          subtapp.GetPodNameCommsBridgeCopy(s.GroupID, subtapp.GetRobotID(i+1)),
				Namespace:                     ns,
				Labels:                        subtapp.GetPodLabelsCommsBridgeCopy(s.GroupID, s.ParentGroupID, r).Map(),
				RestartPolicy:                 orchestrator.RestartPolicyNever,
				TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
				NodeSelector:                  subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
				InitContainers: []orchestrator.Container{
					{
						Name:    "chown-shared-volume",
						Image:   "infrastructureascode/aws-cli:latest",
						Command: []string{"/bin/sh"},
						Args:    []string{"-c", fmt.Sprintf("chown %d:%d /tmp", 1000, 1000)},
						Volumes: volumes,
					},
				},
				Containers: []orchestrator.Container{
					{
						Name:    subtapp.GetContainerNameCommsBridgeCopy(),
						Image:   "infrastructureascode/aws-cli:latest",
						Command: []string{"tail", "-f", "/dev/null"},
						Volumes: volumes,
						EnvVars: subtapp.GetEnvVarsCommsBridgeCopy(
							s.Platform().Store().Ignition().Region(),
							accessKey,
							secretAccessKey,
						),
					},
				},
				Volumes:     volumes,
				Nameservers: s.Platform().Store().Orchestrator().Nameservers(),
			})
		}
	}

	return jobs.LaunchPodsInput(pods), nil
}
