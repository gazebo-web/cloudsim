package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitSimulationPods is a job extending the generic jobs.Wait to wait for all simulation pods to be ready.
var WaitSimulationPods = jobs.Wait.Extend(actions.Job{
	Name:       "wait-simulation-pods",
	PreHooks:   []actions.JobFunc{setStartState, createSimulationWaitRequest},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createSimulationWaitRequest is a pre-hook of the specific WaitSimulationPods job in charge of creating the request for the jobs.Wait job.
func createSimulationWaitRequest(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	res := resource.NewResource("", "", subtapp.GetPodLabelsBase(s.GroupID, s.ParentGroupID))

	// Create wait for condition request
	req := s.Platform().Orchestrator().Pods().WaitForCondition(res, resource.ReadyCondition)

	// Get timeout and poll frequency from store
	timeout := s.Platform().Store().Orchestrator().Timeout()
	pollFreq := s.Platform().Store().Orchestrator().PollFrequency()

	// Return new wait input
	return jobs.WaitInput{
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}
