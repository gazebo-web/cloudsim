package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/gloo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"time"
)

var WaitUpstream = jobs.Wait.Extend(actions.Job{
	Name:            "wait-upstream-gloo",
	PreHooks:        []actions.JobFunc{setStartState, createWaitRequestForUpstream},
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: nil,
	InputType:       nil,
	OutputType:      nil,
})

func createWaitRequestForUpstream(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
	ext := s.Platform().Orchestrator().Extensions().(*gloo.Gloo)
	ns := s.Platform().Store().Orchestrator().IngressNamespace()

	req := waiter.NewWaitRequest(func() (bool, error) {
		res, err := ext.GetUpstream(ns, s.ServiceLabels)
		if err != nil {
			return false, err
		}
		s.UpstreamName = res.Namespace()
		return true, nil
	})

	return jobs.WaitInput{
		Request:       req,
		PollFrequency: time.Second,
		Timeout:       time.Minute,
	}, nil
}
