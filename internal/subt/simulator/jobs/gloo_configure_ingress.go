package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/gloo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

var ConfigureIngressGloo = jobs.ConfigureIngress.Extend(actions.Job{
	Name:            "configure-ingress-gloo",
	PreHooks:        []actions.JobFunc{setStartState, prepareConfigureIngressInputUsingGloo},
	PostHooks:       []actions.JobFunc{checkConfigureIngressError, returnState},
	RollbackHandler: nil,
	InputType:       nil,
	OutputType:      nil,
})

func prepareConfigureIngressInputUsingGloo(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
	name := s.Platform().Store().Orchestrator().IngressName()
	ns := s.Platform().Store().Orchestrator().IngressNamespace()
	host := s.Platform().Store().Orchestrator().IngressHost()

	s.Platform().Orchestrator().Ingresses().GetDestination()

	matcher := gloo.GenerateMatcher("")
	action := gloo.GenerateRouteAction(ns, upstream)
	paths := []orchestrator.Path{gloo.NewPath("", matcher, action)}

	return jobs.ConfigureIngressInput{
		Name:      name,
		Namespace: ns,
		Host:      host,
		Paths:     paths,
	}, nil
}

func checkConfigureIngressError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.ConfigureIngressOutput)
	if out.Error != nil {
		return nil, out.Error
	}
	return nil, nil
}
