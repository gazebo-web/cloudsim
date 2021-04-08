package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/gloo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemoveIngressRulesGloo is a job extending the generic jobs.RemoveIngressRules job that will remove rules from
// the Gloo ingress.
var RemoveIngressRulesGloo = jobs.RemoveIngressRules.Extend(actions.Job{
	Name:       "remove-ingress-rules-gloo",
	PreHooks:   []actions.JobFunc{setStopState, getUpstream, prepareRemoveIngressRulesInput},
	PostHooks:  []actions.JobFunc{checkRemoveIngressRulesError, returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
})

// getUpstream is a prehook in charge of getting the upstream name before removing it.
func getUpstream(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	vs := s.Platform().Orchestrator().Ingresses().(*gloo.VirtualServices)

	ns := s.Platform().Store().Orchestrator().IngressNamespace()

	res, err := vs.GetUpstream(ns, application.GetWebsocketServiceLabels(s.GroupID))
	if err != nil {
		return nil, err
	}

	s.UpstreamName = res.Name()
	store.SetState(s)
	return s, nil
}

// prepareRemoveIngressRulesInput is a pre-hook for the RemoveIngressRulesGloo job in charge of configuring the
// the input for the generic jobs.RemoveIngressRules job.
func prepareRemoveIngressRulesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	name := s.Platform().Store().Orchestrator().IngressName()
	host := s.Platform().Store().Orchestrator().IngressHost()
	ns := s.Platform().Store().Orchestrator().IngressNamespace()

	matcher := gloo.GenerateRegexMatcher(application.GetSimulationIngressPath(s.GroupID))
	action := gloo.GenerateRouteAction(ns, s.UpstreamName)
	paths := []ingresses.Path{gloo.NewPath(s.GroupID.String(), matcher, action)}

	return jobs.RemoveIngressRulesInput{
		Name:      name,
		Namespace: ns,
		Host:      host,
		Paths:     paths,
	}, nil
}

// checkRemoveIngressRulesError checks if the given output from the generic jobs.RemoveIngressRules job returns an error.
func checkRemoveIngressRulesError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.RemoveIngressRulesOutput)
	if out.Error != nil {
		return nil, out.Error
	}
	return nil, nil
}
