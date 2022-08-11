package jobs

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// CreateNetworkPoliciesInput is the input for the CreateNetworkPolicies job.
type CreateNetworkPoliciesInput []network.CreateNetworkPolicyInput

// CreateNetworkPoliciesOutput is the output of the CreateNetworkPolicies job.
type CreateNetworkPoliciesOutput struct {
	// Resource is the representation of the network policies that were created.
	Resource []resource.Resource

	// Error has a reference to the latest error thrown when creating the network policies.
	Error error
}

// CreateNetworkPolicies is a generic job to be used to create network policies.
var CreateNetworkPolicies = &actions.Job{
	Execute: createNetworkPolicies,
}

// createNetworkPolicies is used by the CreateNetworkPolicies job as the execute function.
func createNetworkPolicies(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	input := value.(CreateNetworkPoliciesInput)

	resources := make([]resource.Resource, 0, len(input))
	for _, in := range input {
		res, err := s.Platform().Orchestrator().NetworkPolicies().Create(context.TODO(), in)

		if err != nil {
			return CreateNetworkPoliciesOutput{
				Resource: resources,
				Error:    err,
			}, nil
		}

		resources = append(resources, res)
	}

	return CreateNetworkPoliciesOutput{
		Resource: resources,
		Error:    nil,
	}, nil
}
