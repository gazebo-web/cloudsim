package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// RemoveNetworkPoliciesInput is the input for the RemoveNetworkPolicies job.
type RemoveNetworkPoliciesInput []resource.Resource

// RemoveNetworkPoliciesOutput is the output of the RemoveNetworkPolicies job.
type RemoveNetworkPoliciesOutput struct {
	// Resource is the representation of the network policies that were removed.
	Resource []resource.Resource

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

	resources := make([]resource.Resource, 0, len(input))
	for _, in := range input {
		err := s.Platform().Orchestrator().NetworkPolicies().Remove(in.Name(), in.Namespace())

		if err != nil {
			return RemoveNetworkPoliciesOutput{
				Resource: resources,
				Error:    err,
			}, nil
		}

		resources = append(resources, in)
	}

	return RemoveNetworkPoliciesOutput{
		Resource: resources,
		Error:    nil,
	}, nil
}
