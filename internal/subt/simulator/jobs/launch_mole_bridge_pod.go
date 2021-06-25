package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"strconv"
)

// LaunchMoleBridgePod launches a mole bridge pod.
var LaunchMoleBridgePod = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-mole-bridge-pod",
	PreHooks:        []actions.JobFunc{setStartState, prepareMoleBridgePodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackLaunchMoleBridgePod,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func rollbackLaunchMoleBridgePod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{},
	err error) (interface{}, error) {

	s := store.State().(*state.StartSimulation)

	name := subtapp.GetPodNameMoleBridge(s.GroupID)
	ns := s.Platform().Store().Orchestrator().Namespace()
	_, _ = s.Platform().Orchestrator().Pods().Delete(resource.NewResource(name, ns, nil))

	return nil, nil
}

// prepareMoleBridgePodInput prepares the input for the generic LaunchPods job to launch comms bridge pods.
func prepareMoleBridgePodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	s := store.State().(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	subtSim := sim.(subt.Simulation)
	worldIndex := subtSim.GetWorldIndex()
	runIndex := subtSim.GetRunIndex()

	track, err := s.SubTServices().Tracks().Get(subtSim.GetTrack(), worldIndex, runIndex)
	if err != nil {
		return nil, err
	}

	// Only launch the Mole bridge if the image was defined
	if track.MoleBridgeImage == nil {
		return jobs.LaunchPodsInput([]pods.CreatePodInput{}), nil
	}

	// Prepare the mole pod input
	privileged := true
	allowPrivilegesEscalation := true

	teamID := 0
	owner := subtSim.GetOwner()
	if owner != nil {
		org, em := s.Services().Users().GetOrganization(*owner)
		if em != nil {
			return nil, em.BaseError
		}
		teamID = int(org.ID)
	}

	envVars := map[string]string{
		"PYTHONUNBUFFERED":       "0",
		"CS_PB_PULSAR_IP":        s.Platform().Store().Mole().BridgePulsarAddress(),
		"CS_PB_PULSAR_PORT":      strconv.Itoa(s.Platform().Store().Mole().BridgePulsarPort()),
		"CS_PB_PULSAR_HTTP_PORT": strconv.Itoa(s.Platform().Store().Mole().BridgePulsarHTTPPort()),
		"CS_PB_TOPIC_REGEX":      s.Platform().Store().Mole().BridgeTopicRegex(),
		"ROS_MASTER_URI":         fmt.Sprintf("http://%s:11311", s.GazeboServerIP),
		"WORLD_ID":               strconv.Itoa(worldIndex + 1),
		"TEAM_ID":                strconv.Itoa(teamID),
		"REPLICATION_ID":         string(s.GroupID),
	}

	podInputs := []pods.CreatePodInput{
		{
			Name:                          subtapp.GetPodNameMoleBridge(s.GroupID),
			Namespace:                     s.Platform().Store().Orchestrator().Namespace(),
			Labels:                        subtapp.GetPodLabelsMoleBridge(s.GroupID, s.ParentGroupID).Map(),
			RestartPolicy:                 pods.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsGazeboServer(s.GroupID),
			Containers: []pods.Container{
				{
					Name:                     subtapp.GetContainerNameMoleBridge(),
					Image:                    *track.MoleBridgeImage,
					Privileged:               &privileged,
					AllowPrivilegeEscalation: &allowPrivilegesEscalation,
					EnvVars:                  envVars,
				},
			},
			Nameservers: s.Platform().Store().Orchestrator().Nameservers(),
		},
	}

	return jobs.LaunchPodsInput(podInputs), nil
}
