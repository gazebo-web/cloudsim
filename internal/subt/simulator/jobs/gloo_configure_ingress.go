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

// ConfigureIngressGloo is a job extending the generic jobs.ConfigureIngress job that will configure Gloo to accept
// websocket connections to the gzserver instance.
var ConfigureIngressGloo = jobs.ConfigureIngress.Extend(actions.Job{
	Name:            "configure-ingress-gloo",
	PreHooks:        []actions.JobFunc{setStartState, prepareConfigureIngressInputUsingGloo},
	PostHooks:       []actions.JobFunc{checkConfigureIngressError, returnState},
	RollbackHandler: rollbackGlooIngress,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

// prepareConfigureIngressInputUsingGloo is a pre-hook for the ConfigureIngressGloo job in charge of configuring the
// the input for the generic jobs.ConfigureIngress job.
func prepareConfigureIngressInputUsingGloo(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := s.Platform().Store().Orchestrator().IngressName()
	host := s.Platform().Store().Orchestrator().IngressHost()

	ns := s.Platform().Store().Orchestrator().IngressNamespace()

	matcher := gloo.GenerateRegexMatcher(application.GetSimulationIngressPath(s.GroupID))
	action := gloo.GenerateRouteAction(ns, s.UpstreamName)
	paths := []ingresses.Path{gloo.NewPath(s.GroupID.String(), matcher, action)}

	return jobs.ConfigureIngressInput{
		Name:      name,
		Namespace: ns,
		Host:      host,
		Paths:     paths,
	}, nil
}

// checkConfigureIngressError checks if the given output from the generic ConfigureIngress job returns an error.
func checkConfigureIngressError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(jobs.ConfigureIngressOutput)
	if out.Error != nil {
		return nil, out.Error
	}
	return nil, nil
}

// rollbackGlooIngress is in charge of removing any ingress configuration when there is an error.
func rollbackGlooIngress(store actions.Store, tx *gorm.DB, dep *actions.Deployment, v interface{}, thrownError error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := s.Platform().Store().Orchestrator().IngressName()
	host := s.Platform().Store().Orchestrator().IngressHost()

	ns := s.Platform().Store().Orchestrator().IngressNamespace()

	resource, err := s.Platform().Orchestrator().Ingresses().Get(name, ns)
	if err != nil {
		return nil, nil
	}

	rule, err := s.Platform().Orchestrator().IngressRules().Get(resource, host)
	if err != nil {
		return nil, nil
	}

	matcher := gloo.GenerateRegexMatcher(application.GetSimulationIngressPath(s.GroupID))
	action := gloo.GenerateRouteAction(ns, s.UpstreamName)

	_ = s.Platform().Orchestrator().IngressRules().Remove(rule, gloo.NewPath(s.GroupID.String(), matcher, action))

	return nil, nil
}
