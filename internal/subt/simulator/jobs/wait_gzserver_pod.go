package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForGazeboServerPod waits for the simulation Gazebo pod to finish launching.
var WaitForGazeboServerPod = jobs.Wait.Extend(actions.Job{
	Name:       "wait-gazebo-server-pod",
	PreHooks:   []actions.JobFunc{createWaitRequestForGzServerPod},
	PostHooks:  []actions.JobFunc{checkWaitError, returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createWaitRequestForGzServerPod is the pre hook in charge of passing the needed input to the Wait job.
func createWaitRequestForGzServerPod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	// Create wait for condition request
	req := s.Platform().Orchestrator().Pods().WaitForCondition(s.GazeboServerPod, orchestrator.HasIPStatusCondition)

	// Get timeout and poll frequency from store
	timeout := s.Platform().Store().Machines().Timeout()
	pollFreq := s.Platform().Store().Machines().PollFrequency()

	// Return new wait input
	return jobs.WaitInput{
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}
