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
)

// LaunchCommsBridgeCopyPods launches the list of comms bridge copy pods.
var LaunchCommsBridgeCopyPods = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-comms-bridge-pods",
	PreHooks:        []actions.JobFunc{setStartState, prepareCommsBridgeCreateCopyPodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackPodCreation,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// prepareCommsBridgeCreateCopyPodInput prepares the input for the generic LaunchPods job to launch comms bridge pods.
func prepareCommsBridgeCreateCopyPodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	if !s.Platform().Store().Ignition().LogsCopyEnabled() {
		return jobs.LaunchPodsInput{}, nil
	}

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	subtSim := sim.(subt.Simulation)

	var pods []orchestrator.CreatePodInput

	for i, r := range subtSim.GetRobots() {
		volumes := []orchestrator.Volume{
			{
				Name:         "logs",
				HostPath:     "/tmp",
				SubPath:      "/robot-logs",
				HostPathType: orchestrator.HostPathDirectoryOrCreate,
				MountPath:    s.Platform().Store().Ignition().ROSLogsPath(),
			},
		}

		secretsName := s.Platform().Store().Ignition().SecretsName()
		ns := s.Platform().Store().Orchestrator().Namespace()

		secret, err := s.Platform().Secrets().Get(context.TODO(), secretsName, ns)
		if err != nil {
			return nil, err
		}

		accessKey := string(secret.Data[s.Platform().Store().Ignition().AccessKeyLabel()])
		secretAccessKey := string(secret.Data[s.Platform().Store().Ignition().SecretAccessKeyLabel()])

		pods = append(pods, orchestrator.CreatePodInput{
			Name:                          subtapp.GetPodNameCommsBridgeCopy(s.GroupID, subtapp.GetRobotID(i)),
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

	return jobs.LaunchPodsInput(pods), nil
}
