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
	"time"
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
		pods = append(pods, prepareCreatePodInput(configPod{
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
			envVars: map[string]string{
				"IGN_PARTITION":  s.GroupID.String(),
				"IGN_RELAY":      s.GazeboServerIP,
				"IGN_VERBOSE":    s.Platform().Store().Ignition().Verbosity(),
				"ROBOT_NAME":     r.Name(),
				"IGN_IP":         "", // To be removed.
				"ROS_MASTER_URI": "http://($ROS_IP):11311",
			},
			nameservers: s.Platform().Store().Orchestrator().Nameservers(),
		}))

		if s.Platform().Store().Ignition().LogsCopyEnabled() {
			secretsName := s.Platform().Store().Ignition().SecretsName()
			ns := s.Platform().Store().Orchestrator().Namespace()

			secret, err := s.Platform().Secrets().Get(context.TODO(), secretsName, ns)
			if err != nil {
				return nil, err
			}

			accessKey := secret.Data[s.Platform().Store().Ignition().AccessKeyLabel()]
			secretAccessKey := secret.Data[s.Platform().Store().Ignition().SecretAccessKeyLabel()]

			// Create copy pod input
			pods = append(pods, prepareBridgeCopyCreatePodInput(configBridgeCopyPod{
				name:                   subtapp.GetPodNameCommsBridgeCopy(s.GroupID, subtapp.GetRobotID(i+1)),
				namespace:              s.Platform().Store().Orchestrator().Namespace(),
				labels:                 subtapp.GetPodLabelsCommsBridgeCopy(s.GroupID, s.ParentGroupID, r).Map(),
				terminationGracePeriod: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
				nodeSelector:           subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
				hostLogsPath:           logMountPath,
				mountLogsPath:          s.Platform().Store().Ignition().ROSLogsPath(),
				region:                 s.Platform().Store().Ignition().Region(),
				accessKey:              string(accessKey),
				secretAccessKey:        string(secretAccessKey),
				nameservers:            s.Platform().Store().Orchestrator().Nameservers(),
			}))
		}
	}

	return pods, nil
}

type configBridgeCopyPod struct {
	name                   string
	namespace              string
	labels                 map[string]string
	terminationGracePeriod time.Duration
	nodeSelector           orchestrator.Selector
	hostLogsPath           string
	mountLogsPath          string
	region                 string
	accessKey              string
	secretAccessKey        string
	nameservers            []string
}

func prepareBridgeCopyCreatePodInput(c configBridgeCopyPod) orchestrator.CreatePodInput {
	return prepareCreatePodInput(configPod{
		name:                   c.name,
		namespace:              c.namespace,
		labels:                 c.labels,
		restartPolicy:          orchestrator.RestartPolicyNever,
		terminationGracePeriod: c.terminationGracePeriod,
		nodeSelector:           c.nodeSelector,
		containerName:          "copy-to-s3",
		image:                  "infrastructureascode/aws-cli:latest",
		command:                []string{"tail", "-f", "/dev/null"},
		volumes: []orchestrator.Volume{
			{
				Name:         "logs",
				HostPath:     c.hostLogsPath,
				MountPath:    c.mountLogsPath,
				HostPathType: orchestrator.HostPathDirectoryOrCreate,
			},
		},
		envVars: map[string]string{
			"AWS_DEFAULT_REGION":    c.region,
			"AWS_REGION":            c.region,
			"AWS_ACCESS_KEY_ID":     c.accessKey,
			"AWS_SECRET_ACCESS_KEY": c.secretAccessKey,
		},
		nameservers: c.nameservers,
	})
}
