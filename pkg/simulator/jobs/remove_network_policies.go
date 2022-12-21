package jobs

import (
	"context"
	"github.com/gazebo-web/cloudsim/pkg/actions"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/resource"
	"github.com/gazebo-web/cloudsim/pkg/simulator/state"
	"github.com/jinzhu/gorm"
)

// RemoveNetworkPoliciesInput is the input for the RemoveNetworkPolicies job.
type RemoveNetworkPoliciesInput struct {
	Selector  resource.Selector
	Namespace string
}

// RemoveNetworkPoliciesOutput is the output of the RemoveNetworkPolicies job.
type RemoveNetworkPoliciesOutput struct {
	// Error has a reference to the latest error thrown when removing the network policies.
	Error error
}

// RemoveNetworkPolicies is a generic job to be used to remove network policies.
var RemoveNetworkPolicies = &actions.Job{
	Execute: removeNetworkPolicies,
}

// removeNetworkPolicies is used by the RemoveNetworkPolicies job as the execute function.
func removeNetworkPolicies(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	input := value.(RemoveNetworkPoliciesInput)

	err := s.Platform().Orchestrator().NetworkPolicies().RemoveBulk(context.TODO(), input.Namespace, input.Selector)

	return RemoveNetworkPoliciesOutput{
		Error: err,
	}, nil
}
