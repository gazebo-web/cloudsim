package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitRobots is a job extending the generic jobs.Wait to wait for robots to be ready.
var WaitRobots = jobs.Wait.Extend(actions.Job{
	Name:       "wait-robots",
	PreHooks:   []actions.JobFunc{setStartState, createWaitRequestForRobots},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createWaitRequestForRobots is a pre-hook of the specific WaitRobots job in charge of creating the request for the jobs.Wait job.
func createWaitRequestForRobots(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	res := orchestrator.NewResource("", "", subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID))

	// Create wait for condition request
	req := s.Platform().Orchestrator().Pods().WaitForCondition(res, orchestrator.ReadyCondition)

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
