package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
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

	return jobs.RemoveNetworkPoliciesInput{
		Namespace: ns,
		Selector:  subtapp.GetPodLabelsBase(s.GroupID, nil),
	}, nil
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
