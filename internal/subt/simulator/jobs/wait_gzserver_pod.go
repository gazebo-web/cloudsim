package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

var WaitForGazeboServerPod = Wait.Extend(actions.Job{
	Name:     "wait-gazebo-server-pod",
	PreHooks: []actions.JobFunc{createWaitRequestForGzServerPod},
})

func createWaitRequestForGzServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	simCtx := context.NewContext(ctx)

	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	labels := map[string]string{
		"cloudsim":          "true",
		"cloudsim-group-id": string(gid),
		"gzserver":          "true",
	}

	namespace := simCtx.Platform().Store().Orchestrator().Namespace()

	res := orchestrator.NewResource("", namespace, orchestrator.NewSelector(labels))

	req := simCtx.Platform().Orchestrator().Pods().WaitForCondition(res, orchestrator.ReadyCondition)

	timeout := simCtx.Platform().Store().Machines().Timeout()
	pollFreq := simCtx.Platform().Store().Machines().PollFrequency()

	return WaitInput{
		GroupID:       gid,
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}
