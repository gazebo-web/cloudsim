package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForGazeboServerPod extends the Wait job to fill the input data needed by Wait's execute function.
var WaitForGazeboServerPod = jobs.Wait.Extend(actions.Job{
	Name:       "wait-gazebo-server-pod",
	PreHooks:   []actions.JobFunc{createWaitRequestForGzServerPod},
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
})

// createWaitRequestForGzServerPod is the prehook in charge of passing the needed input to the Wait job.
func createWaitRequestForGzServerPod(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Get context
	simCtx := context.NewContext(ctx)

	data := value.(*StartSimulationData)

	// Create wait for condition request
	req := simCtx.Platform().Orchestrator().Pods().WaitForCondition(data.GazeboServerPod, orchestrator.HasIPStatusCondition)

	// Get timeout and poll frequency from store
	timeout := simCtx.Platform().Store().Machines().Timeout()
	pollFreq := simCtx.Platform().Store().Machines().PollFrequency()

	// Return new wait input
	return jobs.WaitInput{
		GroupID:       data.GroupID,
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}
