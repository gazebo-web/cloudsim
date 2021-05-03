package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/gloo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"time"
)

// WaitUpstream is a job extending the generic jobs.Wait to wait for an upstream to be available.
var WaitUpstream = jobs.Wait.Extend(actions.Job{
	Name:       "wait-upstream-gloo",
	PreHooks:   []actions.JobFunc{setStartState, createWaitRequestForUpstream},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createWaitRequestForUpstream is a pre-hook of the specific WaitUpstream job in charge of creating the request for the jobs.Wait job.
func createWaitRequestForUpstream(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	vs := s.Platform().Orchestrator().Ingresses().(*gloo.VirtualServices)

	ns := s.Platform().Store().Orchestrator().IngressNamespace()

	req := waiter.NewWaitRequest(func() (bool, error) {
		// The upstream sometimes takes a while to be ready. If this fails to get the upstream, return false and NOT
		// an error as that would terminate the waiter. This waiter will only fail if it times out.
		res, err := vs.GetUpstream(ns, application.GetWebsocketServiceLabels(s.GroupID))
		if err != nil {
			return false, nil
		}
		s.UpstreamName = res.Name()

		// Update the state
		store.SetState(s)

		return true, nil
	})

	return jobs.WaitInput{
		Request:       req,
		PollFrequency: 10 * time.Second,
		Timeout:       10 * time.Minute,
	}, nil
}
