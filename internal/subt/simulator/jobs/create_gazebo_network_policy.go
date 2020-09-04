package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// CreateGazeboServerNetworkPolicy is used to create a network policy in the orchestrator.Cluster for the gazebo server.
var CreateGazeboServerNetworkPolicy = &actions.Job{
	Name:            "create-gazebo-network-policy",
	PreHooks:        []actions.JobFunc{prepareGazeboNetworkPolicyInput},
	Execute:         createGazeboNetworkPolicy,
	RollbackHandler: rollbackCreateGazeboNetworkPolicy,
	InputType:       actions.GetJobDataType(simulations.GroupID("")),
	OutputType:      actions.GetJobDataType(simulations.GroupID("")),
}

func prepareGazeboNetworkPolicyInput(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	simCtx := context.NewContext(ctx)

	sim, err := simCtx.Services().Simulations().Get(gid)
	if err != nil {
		return nil, err
	}

	// Set up pod name
	podName := "prefix-groupid-gzserver"

	// Set up labels
	baseLabels := map[string]string{
		"cloudsim":          "true",
		"pod-group":         podName,
		"cloudsim-group-id": string(gid),
		"SubT":              "true",
	}

	gzServerLabels := map[string]string{
		"cloudsim":          "true",
		"pod-group":         podName,
		"cloudsim-group-id": string(gid),
		"SubT":              "true",
		"gzserver":          "true",
	}
	podSelector := orchestrator.NewSelector(gzServerLabels)

	bridgeLabels := map[string]string{
		"cloudsim":          "true",
		"pod-group":         podName,
		"cloudsim-group-id": string(gid),
		"SubT":              "true",
		"comms-bridge":      "true",
	}

	cidr := fmt.Sprintf("%s/32", simCtx.Platform().Store().Ignition().IP())

	input := orchestrator.CreateNetworkPolicyInput{
		Name:        fmt.Sprintf("%s-%s-%s", "network-policy", sim.GroupID(), "gzserver"),
		Namespace:   simCtx.Platform().Store().Orchestrator().Namespace(),
		Labels:      baseLabels,
		PodSelector: podSelector,
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
			orchestrator.NewSelector(bridgeLabels),
		},
		PeersTo: []orchestrator.Selector{
			orchestrator.NewSelector(bridgeLabels),
		},
	}

	return map[string]interface{}{
		"groupID":                  gid,
		"createNetworkPolicyInput": input,
	}, nil
}

func createGazeboNetworkPolicy(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Parse input
	inputMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	createPolicyInput, ok := inputMap["createNetworkPolicyInput"].(orchestrator.CreateNetworkPolicyInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	gid, ok := inputMap["groupID"].(simulations.GroupID)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	simCtx := context.NewContext(ctx)

	_, err := simCtx.Platform().Orchestrator().NetworkPolicies().Create(createPolicyInput)
	if err != nil {

		return nil, err
	}
	return gid, nil
}

func rollbackCreateGazeboNetworkPolicy(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	// TODO: Remove network policy.
	return value, nil
}
