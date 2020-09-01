package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// CreateGazeboServerNetworkPolicy is used to create a network policy in the orchestrator.Cluster for the gazebo server.
var CreateGazeboServerNetworkPolicy = &actions.Job{
	Name:            "create-gazebo-network-policy",
	Execute:         createGazeboNetworkPolicy,
	RollbackHandler: rollbackCreateGazeboNetworkPolicy,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

func createGazeboNetworkPolicy(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	input := orchestrator.CreateNetworkPolicyInput{
		Name: fmt.Sprintf("%s-%s-%s", "network-policy", sim.GroupID(), "gzserver"),
		// Namespace:     simCtx.Platform().Store().Cluster().Namespace(),
		Labels:      nil,
		PodSelector: nil,
		// CIDR:          simCtx.Platform().Store().Ignition().IP(),
		WebsocketPort: 0,
		PeersFrom:     nil,
		PeersTo:       nil,
	}

	_, err = simCtx.Platform().Orchestrator().NetworkPolicies().Create(input)
	if err != nil {

		return nil, err
	}
	return gid, nil
}

func rollbackCreateGazeboNetworkPolicy(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	// TODO: Remove network policy.
	return value, nil
}
