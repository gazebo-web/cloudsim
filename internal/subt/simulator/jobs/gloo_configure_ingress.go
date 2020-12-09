package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/gloo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// ConfigureIngressGloo is a job extending the generic jobs.ConfigureIngress job that will configure Gloo to accept
// websocket connections to the gzserver instance.
var ConfigureIngressGloo = jobs.ConfigureIngress.Extend(actions.Job{
	Name:            "configure-ingress-gloo",
	PreHooks:        []actions.JobFunc{setStartState, prepareConfigureIngressInputUsingGloo},
	PostHooks:       []actions.JobFunc{checkConfigureIngressError, returnState},
	RollbackHandler: nil,
	InputType:       nil,
	OutputType:      nil,
})

// prepareConfigureIngressInputUsingGloo is a pre-hook for the ConfigureIngressGloo job in charge of configuring the
// the input for the generic jobs.ConfigureIngress job.
func prepareConfigureIngressInputUsingGloo(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := s.Platform().Store().Orchestrator().IngressName()
	host := s.Platform().Store().Orchestrator().IngressHost()

	ns := s.Platform().Store().Orchestrator().Namespace()

	matcher := gloo.GenerateRegexMatcher("")
	action := gloo.GenerateRouteAction(ns, s.UpstreamName)
	paths := []orchestrator.Path{gloo.NewPath("", matcher, action)}

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
