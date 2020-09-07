package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// WaitForGazeboServerPod extends the Wait job to fill the input data needed by Wait's execute function.
var WaitForGazeboServerPod = Wait.Extend(actions.Job{
	Name:       "wait-gazebo-server-pod",
	PreHooks:   []actions.JobFunc{createWaitRequestForGzServerPod},
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
})

// createWaitRequestForGzServerPod is the prehook in charge of passing the needed input to the Wait job.
func createWaitRequestForGzServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Parse group id
	gid, ok := value.(simulations.GroupID)
	if !ok {
		return nil, simulations.ErrInvalidGroupID
	}

	// Generate gzserver labels
	labels := map[string]string{
		"cloudsim":          "true",
		"cloudsim-group-id": string(gid),
		"gzserver":          "true",
	}

	// Get context
	simCtx := context.NewContext(ctx)

	// Get default namespace
	namespace := simCtx.Platform().Store().Orchestrator().Namespace()

	// Create resource
	// TODO: Add name
	res := orchestrator.NewResource("", namespace, orchestrator.NewSelector(labels))

	// Create wait for condition request
	req := simCtx.Platform().Orchestrator().Pods().WaitForCondition(res, orchestrator.ReadyCondition)

	// Get timeout and poll frequency from store
	timeout := simCtx.Platform().Store().Machines().Timeout()
	pollFreq := simCtx.Platform().Store().Machines().PollFrequency()

	// Return new wait input
	return WaitInput{
		GroupID:       gid,
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}
