package jobs

import (
	"context"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// LaunchMappingServerCopyPod launches a mapping server copy pod.
var LaunchMappingServerCopyPod = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-mapping-copy-pod",
	PreHooks:        []actions.JobFunc{setStartState, prepareMappingCopyPodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackLaunchMappingServerCopyPod,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func rollbackLaunchMappingServerCopyPod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := subtapp.GetPodNameMappingServerCopy(s.GroupID)
	ns := s.Platform().Store().Orchestrator().Namespace()

	_, _ = s.Platform().Orchestrator().Pods().Delete(resource.NewResource(name, ns, nil))

	return nil, nil
}

func prepareMappingCopyPodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	if !s.Platform().Store().Ignition().LogsCopyEnabled() {
		return jobs.LaunchPodsInput{}, nil
	}

	if !isMappingServerEnabled(s.SubTServices(), s.GroupID) {
		return jobs.LaunchPodsInput{}, nil
	}

	// Set up namespace
	namespace := s.Platform().Store().Orchestrator().Namespace()

	// Set up nameservers
	nameservers := s.Platform().Store().Orchestrator().Nameservers()

	// Set up secrets
	secretsName := s.Platform().Store().Ignition().SecretsName()
	secret, err := s.Platform().Secrets().Get(context.TODO(), secretsName, namespace)
	if err != nil {
		return nil, err
	}

	accessKey := string(secret.Data[s.Platform().Store().Ignition().AccessKeyLabel()])
	secretAccessKey := string(secret.Data[s.Platform().Store().Ignition().SecretAccessKeyLabel()])

	volumes := []pods.Volume{
		{
			Name:         "logs",
			HostPath:     "/tmp/mapping",
			HostPathType: pods.HostPathDirectoryOrCreate,
			MountPath:    "/tmp/mapping",
		},
	}

	return jobs.LaunchPodsInput{
		{
			Name:                          subtapp.GetPodNameMappingServerCopy(s.GroupID),
			Namespace:                     namespace,
			Labels:                        subtapp.GetPodLabelsMappingServerCopy(s.GroupID, s.ParentGroupID).Map(),
			RestartPolicy:                 pods.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsGazeboServer(s.GroupID),
			Containers: []pods.Container{
				{
					Name:    subtapp.GetContainerNameMappingServerCopy(),
					Image:   "infrastructureascode/aws-cli:latest",
					Command: []string{"tail", "-f", "/dev/null"},
					Volumes: volumes,
					EnvVars: subtapp.GetEnvVarsMappingServerCopy(
						s.Platform().Store().Ignition().Region(),
						accessKey,
						secretAccessKey,
					),
				},
			},
			Volumes:     volumes,
			Nameservers: nameservers,
		},
	}, nil
}
