package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemoveNetworkPolicies is a specific subt job to remove network policies applied to the ignition gazebo server,
// field computer pods and comms bridges.
var RemoveNetworkPolicies = jobs.RemoveNetworkPolicies.Extend(actions.Job{
	Name:       "remove-network-policies",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemoveNetworkPoliciesInput},
	PostHooks:  []actions.JobFunc{checkRemoveNetworkPoliciesError, returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
})

// prepareRemoveNetworkPoliciesInput prepares the input for the generic jobs.RemoveNetworkPolicies job.
// It's a pre-hook of the RemoveNetworkPolicies job.
func prepareRemoveNetworkPoliciesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	ns := s.Platform().Store().Orchestrator().Namespace()

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	// This job is removing the following network policies:
	// 1 Network policy for the Ignition Gazebo Server
	// 2 Network policies per robot pod:
	// 		- Field computer network policy
	// 		- Comms bridge network policy
	resources := make([]orchestrator.Resource, 0, 2*len(robots)+1)

	for i := range robots {
		robotID := subtapp.GetRobotID(i)
		resources = append(resources, orchestrator.NewResource(subtapp.GetPodNameCommsBridge(s.GroupID, robotID), ns, nil))
		resources = append(resources, orchestrator.NewResource(subtapp.GetPodNameFieldComputer(s.GroupID, robotID), ns, nil))
	}

	resources = append(resources, orchestrator.NewResource(subtapp.GetPodNameGazeboServer(s.GroupID), ns, nil))

	return jobs.RemoveNetworkPoliciesInput(resources), nil
}

// checkRemoveNetworkPoliciesError checks if an error has been thrown while removing network policies.
// It's a post-hook of the RemoveNetworkPolicies job.
func checkRemoveNetworkPoliciesError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.RemoveNetworkPoliciesOutput)
	if out.Error != nil {
		return nil, out.Error
	}
	return nil, nil
}
