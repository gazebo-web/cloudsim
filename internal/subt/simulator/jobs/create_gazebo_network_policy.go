package jobs

import (
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

	input := orchestrator.CreateNetworkPolicyInput{
		Name:        "",
		Labels:      map[string]string{},
		PodSelector: orchestrator.Selector(),
	}

	simCtx.Platform().Orchestrator().NetworkPolicies().Create(input)

}

func rollbackCreateGazeboNetworkPolicy(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {

}
