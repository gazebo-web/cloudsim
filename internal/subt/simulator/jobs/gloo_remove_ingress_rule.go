package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/gloo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// ConfigureIngressGloo is a job extending the generic jobs.ConfigureIngress job that will configure Gloo to accept
// websocket connections to the gzserver instance.
var RemoveIngressRulesGloo = jobs.RemoveIngressRules.Extend(actions.Job{
	Name:       "remove-ingress-rules-gloo",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemoveIngressRulesInput},
	PostHooks:  []actions.JobFunc{checkRemoveIngressRulesError, returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
})

// prepareRemoveIngressRulesInput is a pre-hook for the RemoveIngressRulesGloo job in charge of configuring the
// the input for the generic jobs.RemoveIngressRules job.
func prepareRemoveIngressRulesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := s.Platform().Store().Orchestrator().IngressName()
	host := s.Platform().Store().Orchestrator().IngressHost()

	ns := s.Platform().Store().Orchestrator().Namespace()

	matcher := gloo.GenerateRegexMatcher(application.GetSimulationIngressPath(s.GroupID))
	action := gloo.GenerateRouteAction(ns, s.UpstreamName)
	paths := []orchestrator.Path{gloo.NewPath(s.GroupID.String(), matcher, action)}

	return jobs.ConfigureIngressInput{
		Name:      name,
		Namespace: ns,
		Host:      host,
		Paths:     paths,
	}, nil
}

// checkRemoveIngressRulesError checks if the given output from the generic jobs.RemoveIngressRules job returns an error.
func checkRemoveIngressRulesError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.ConfigureIngressOutput)
	if out.Error != nil {
		return nil, out.Error
	}
	return nil, nil
}
