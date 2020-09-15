package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// JobCreateNetworkPolicyInput is the input for the CreateNetworkPolicy job.
type JobCreateNetworkPolicyInput orchestrator.CreateNetworkPolicyInput

// JobCreateNetworkPolicyOutput is the output of the CreateNetworkPolicy job.
type JobCreateNetworkPolicyOutput struct {
	// Resource is the representation of the network policy source that was created.
	Resource orchestrator.Resource

	// Error has a reference to the thrown error when creating a network policy.
	Error error
}

// CreateNetworkPolicy is a generic job to be used to create network policies.
var CreateNetworkPolicy = &actions.Job{
	Execute: createNetworkPolicy,
}

// createNetworkPolicy is used by the CreateNetworkPolicy job as the execute function.
func createNetworkPolicy(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.Platform)

	input := value.(JobCreateNetworkPolicyInput)

	createInput := orchestrator.CreateNetworkPolicyInput(input)
	res, err := s.Platform().Orchestrator().NetworkPolicies().Create(createInput)

	return JobCreateNetworkPolicyOutput{
		Resource: res,
		Error:    err,
	}, nil
}
