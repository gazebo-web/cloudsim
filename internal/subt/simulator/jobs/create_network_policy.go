package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// CreateNetworkPolicy is a generic job used to create a network policy in the orchestrator.Cluster.
var CreateNetworkPolicy = &actions.Job{
	Name:    "create-network-policy",
	Execute: createNetworkPolicy,
}

type CreateNetworkPolicyOutput struct {
	Resource orchestrator.Resource
	Error    error
}

func createNetworkPolicy(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	createPolicyInput, ok := value.(orchestrator.CreateNetworkPolicyInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	simCtx := context.NewContext(ctx)

	res, err := simCtx.Platform().Orchestrator().NetworkPolicies().Create(createPolicyInput)

	return CreateNetworkPolicyOutput{
		Resource: res,
		Error:    err,
	}, nil
}
