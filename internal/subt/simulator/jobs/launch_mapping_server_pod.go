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
)

// LaunchMappingServerPod launches a mapping server pod.
var LaunchMappingServerPod = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-mapping-server-pod",
	PreHooks:        []actions.JobFunc{setStartState, prepareMappingCreatePodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackLaunchMappingServerPod,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func rollbackLaunchMappingServerPod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := subtapp.GetPodNameMappingServer(s.GroupID)
	ns := s.Platform().Store().Orchestrator().Namespace()

	_, _ = s.Platform().Orchestrator().Pods().Delete(resource.NewResource(name, ns, nil))

	return nil, nil
}

// prepareMappingCreatePodInput is in charge of preparing the input of jobs.LaunchPods with specific config for launching
// a mapping server pod.
func prepareMappingCreatePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
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
	track, err := s.SubTServices().Tracks().Get(subtSim.GetTrack(), subtSim.GetWorldIndex(), subtSim.GetRunIndex())
	if err != nil {
		return nil, err
	}

	// Generate mapping server command args
	runCommand, err := cmdgen.MapAnalysis(cmdgen.MapAnalysisConfig{
		World:  track.World,
		Robots: subtSim.GetRobots(),
	})
	if err != nil {
		return nil, err
	}

	// Set up container configuration
	privileged := true
	allowPrivilegeEscalation := true

	// TODO: Get ports from Ignition Store
	ports := []int32{11311}

	return jobs.LaunchPodsInput{
		{
			Name:                          subtapp.GetPodNameMappingServer(s.GroupID),
			Namespace:                     namespace,
			Labels:                        subtapp.GetPodLabelsMappingServer(s.GroupID, s.ParentGroupID).Map(),
			RestartPolicy:                 pods.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsGazeboServer(s.GroupID),
			Containers: []pods.Container{
				{
					Name:                     subtapp.GetContainerNameMappingServer(),
					Image:                    "subt/map_analysis:latest",
					Args:                     runCommand,
					Privileged:               &privileged,
					AllowPrivilegeEscalation: &allowPrivilegeEscalation,
					Ports:                    ports,
					EnvVarsFrom:              subtapp.GetEnvVarsFromSourceMappingServer(),
					EnvVars:                  subtapp.GetEnvVarsMappingServer(s.GroupID),
				},
			},
		},
	}, nil
}
