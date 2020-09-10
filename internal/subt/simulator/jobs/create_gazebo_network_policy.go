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
var CreateGazeboServerNetworkPolicy = CreateNetworkPolicy.Extend(actions.Job{
	Name:            "create-gazebo-network-policy",
	PreHooks:        []actions.JobFunc{prepareGazeboNetworkPolicyInput},
	PostHooks:       []actions.JobFunc{createGazeboServerNetworkPolicyPostHook},
	RollbackHandler: rollbackCreateGazeboNetworkPolicy,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
})

func prepareGazeboNetworkPolicyInput(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	data := simCtx.Store().Get().(*StartSimulationData)

	// Get simulation
	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	// Set up pod name
	podName := "%s-%s-gzserver"

	cidr := fmt.Sprintf("%s/32", simCtx.Platform().Store().Ignition().IP())

	input := orchestrator.CreateNetworkPolicyInput{
		Name:        fmt.Sprintf("%s-%s-%s", "network-policy", sim.GroupID(), "gzserver"),
		Namespace:   simCtx.Platform().Store().Orchestrator().Namespace(),
		Labels:      data.BaseLabels,
		PodSelector: data.GazeboPodResource.Selector(),
		Ingresses: orchestrator.NetworkIngressRule{
			Ports:    []int32{9002},
			IPBlocks: []string{cidr},
		},
		Egresses: orchestrator.NetworkEgressRule{
			Ports:         nil,
			IPBlocks:      []string{cidr},
			AllowOutbound: true,
		},
		PeersFrom: []orchestrator.Selector{
			orchestrator.NewSelector(data.BridgeLabels),
		},
		PeersTo: []orchestrator.Selector{
			orchestrator.NewSelector(data.BridgeLabels),
		},
	}

	return input, nil
}

func createGazeboServerNetworkPolicyPostHook(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Parse execute function output
	output := value.(CreateNetworkPolicyOutput)

	// Get data from action store
	data := ctx.Store().Get().(StartSimulationData)

	// Assign output resource to the store
	data.GazeboNetworkPolicyResource = output.Resource

	// Persist data
	err := ctx.Store().Set(data)
	if err != nil {
		return nil, err
	}

	// Check if the execute function returned an error
	if output.Error != nil {
		return nil, err
	}

	return data.GroupID, nil
}

func rollbackCreateGazeboNetworkPolicy(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	simCtx := context.NewContext(ctx)

	// Get data from action store
	data := simCtx.Store().Get().(*StartSimulationData)

	_, delErr := simCtx.Platform().Orchestrator().NetworkPolicies().Delete(data.GazeboNetworkPolicyResource)
	if delErr != nil {
		return nil, delErr
	}

	return data.GroupID, nil
}
