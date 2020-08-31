package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// WaitForKubernetesNodes is used to wait until all required kubernetes nodes are ready.
var WaitForKubernetesNodes = &actions.Job{
	Name:       "wait-for-kubernetes-nodes",
	Execute:    waitForKubernetesNodes,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

// waitForKubernetesNodes is the main process executed by WaitForKubernetesNodes.
func waitForKubernetesNodes(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}) (interface{}, error) {

	simCtx := context.NewContext(ctx)

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	res := orchestrator.NewResource("", "", orchestrator.NewSelector(map[string]string{
		"cloudsim_groupid": string(gid),
	}))

	req := simCtx.Platform().Orchestrator().Nodes().WaitForCondition(res, orchestrator.ReadyCondition)

	timeout := simCtx.Platform().Store().Machines().Timeout()
	pollFreq := simCtx.Platform().Store().Machines().PollFrequency()

	err := req.Wait(timeout, pollFreq)
	if err != nil {
		return nil, err
	}

	return gid, nil
}
