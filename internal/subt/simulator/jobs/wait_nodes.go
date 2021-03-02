package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForNodes is the job in charge of waiting until all the simulation instances have joined the cluster
// and have been marked as ready.
var WaitForNodes = jobs.Wait.Extend(actions.Job{
	Name:       "wait-for-nodes",
	PreHooks:   []actions.JobFunc{setStartState, createWaitForNodesInput},
	PostHooks:  []actions.JobFunc{checkWaitError, returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createWaitForNodesInput creates the input for the main wait job. This JobFunc is a pre-hook of the WaitForNodes job.
func createWaitForNodesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(state.StartSimulation)

	store.SetState(s)

	selector := application.GetNodeLabelsBase(s.GroupID)

	res := resource.NewResource("", "", selector)

	w := s.Platform().Orchestrator().Nodes().WaitForCondition(res, resource.ReadyCondition)

	return &jobs.WaitInput{
		Request:       w,
		PollFrequency: s.Platform().Store().Machines().PollFrequency(),
		Timeout:       s.Platform().Store().Machines().Timeout(),
	}, nil
}
