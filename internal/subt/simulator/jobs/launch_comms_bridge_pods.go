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
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"path"
	"time"
)

// LaunchCommsBridge launches the list of comms bridge pods.
var LaunchCommsBridge = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-comms-bridge-pods",
	PreHooks:        []actions.JobFunc{setStartState, prepareCommsBridgePodInput, prepareFieldComputerPodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackPodsCreation,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func prepareCommsBridgePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {

	s := store.State().(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {

	}

	subtSim := sim.(subt.Simulation)

	track, err := s.Services().Tracks().Get(subtSim.Track())

	var pods []orchestrator.CreatePodInput

	for i, r := range subtSim.Robots() {
		childMarsupial := "false"
		if subt.IsRobotChildMarsupial(subtSim.Marsupials(), r) {
			childMarsupial = "true"
		}

		hostPath := "/tmp"

		logDirectory := "robot-logs"

		logMountPath := path.Join(hostPath, logDirectory)

		// Create comms bridge input
		pods = append(pods, prepareCommsBridgeCreatePodInput(configCommsBridgePod{
			groupID:                s.GroupID,
			robotNumber:            i,
			robotID:                subtapp.GetRobotID(i + 1),
			namespace:              s.Platform().Store().Orchestrator().Namespace(),
			labels:                 subtapp.GetPodLabelsCommsBridge(s.GroupID, s.ParentGroupID, r).Map(),
			terminationGracePeriod: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			nodeSelector:           subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
			containerImage:         track.BridgeImage,
			gzServerPodIP:          s.GazeboServerIP,
			robotName:              r.Name(),
			robotType:              r.Kind(),
			ignVerbose:             s.Platform().Store().Ignition().Verbosity(),
			ignIP:                  "", // To be removed.
			mountLogsPath:          s.Platform().Store().Ignition().ROSLogsPath(),
			hostLogsPath:           logMountPath,
			nameservers:            s.Platform().Store().Orchestrator().Nameservers(),
			worldName:              track.World,
			childMarsupial:         childMarsupial,
		}))

		if s.Platform().Store().Ignition().LogsCopyEnabled() {
			secretsName := s.Platform().Store().Ignition().SecretsName()
			ns := s.Platform().Store().Orchestrator().Namespace()

			secret, err := s.Platform().Secrets().Get(context.TODO(), secretsName, ns)
			if err != nil {
				return nil, err
			}

			accessKey := secret.Data["aws-access-key-id"]
			secretAccessKey := secret.Data["aws-secret-access-key"]

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

type configCommsBridgePod struct {
	groupID simulations.GroupID
	// robotNumber is the robot index that's being created from the list of robots when iterating over the list.
	robotNumber int
	// robotID is the robot ID that will be used in the pod name.
	robotID                string
	namespace              string
	labels                 map[string]string
	terminationGracePeriod time.Duration
	nodeSelector           orchestrator.Selector
	containerImage         string
	gzServerPodIP          string
	robotName              string
	robotType              string
	ignVerbose             string
	ignIP                  string // To be removed.
	mountLogsPath          string
	hostLogsPath           string
	nameservers            []string
	worldName              string
	childMarsupial         string
}

func prepareCommsBridgeCreatePodInput(c configCommsBridgePod) orchestrator.CreatePodInput {
	in := configPod{
		name:                   subtapp.GetPodNameCommsBridge(c.groupID, c.robotID),
		namespace:              c.namespace,
		labels:                 c.labels,
		restartPolicy:          orchestrator.RestartPolicyNever,
		terminationGracePeriod: c.terminationGracePeriod,
		nodeSelector:           c.nodeSelector,
		containerName:          "comms-bridge",
		image:                  c.containerImage,
		args: []string{
			c.worldName,
			fmt.Sprintf("robotName%d:=%s", c.robotNumber, c.robotName),
			fmt.Sprintf("robotConfig%d:=%s", c.robotNumber, c.robotType),
			"headless:=true",
			fmt.Sprintf("marsupial:=%s", c.childMarsupial),
		},
		privileged:                true,
		allowPrivilegesEscalation: true,
		volumes: []orchestrator.Volume{
			{
				Name:         "logs",
				HostPath:     c.hostLogsPath,
				HostPathType: orchestrator.HostPathDirectoryOrCreate,
				MountPath:    c.mountLogsPath,
			},
		},
		envVars: map[string]string{
			"IGN_PARTITION":  c.groupID.String(),
			"IGN_RELAY":      c.gzServerPodIP,
			"IGN_VERBOSE":    c.ignVerbose,
			"ROBOT_NAME":     c.robotName,
			"IGN_IP":         c.ignIP, // To be removed.
			"ROS_MASTER_URI": "http://($ROS_IP):11311",
		},
		nameservers: c.nameservers,
	}

	return preparePod(in)
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
	return preparePod(configPod{
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

type configPod struct {
	name                      string
	namespace                 string
	labels                    map[string]string
	restartPolicy             orchestrator.RestartPolicy
	terminationGracePeriod    time.Duration
	nodeSelector              orchestrator.Selector
	containerName             string
	image                     string
	command                   []string
	args                      []string
	privileged                bool
	allowPrivilegesEscalation bool
	ports                     []int32
	volumes                   []orchestrator.Volume
	envVars                   map[string]string
	nameservers               []string
}

func preparePod(c configPod) orchestrator.CreatePodInput {
	return orchestrator.CreatePodInput{
		Name:                          c.name,
		Namespace:                     c.namespace,
		Labels:                        c.labels,
		RestartPolicy:                 c.restartPolicy,
		TerminationGracePeriodSeconds: c.terminationGracePeriod,
		NodeSelector:                  c.nodeSelector,
		Containers: []orchestrator.Container{
			{
				Name:                     c.containerName,
				Image:                    c.image,
				Command:                  c.command,
				Args:                     c.args,
				Privileged:               &c.privileged,
				AllowPrivilegeEscalation: &c.allowPrivilegesEscalation,
				Ports:                    c.ports,
				Volumes:                  c.volumes,
				EnvVars:                  c.envVars,
			},
		},
		Volumes:     c.volumes,
		Nameservers: c.nameservers,
	}
}
