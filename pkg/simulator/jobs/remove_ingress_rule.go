package jobs

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
	"time"
)

// RemoveIngressRulesInput is the input for the RemoveIngressRule job.
type RemoveIngressRulesInput struct {
	Name      string
	Namespace string
	Host      string
	Paths     []ingresses.Path
}

// RemoveIngressRulesOutput is the output of the ConfigureIngress job.
type RemoveIngressRulesOutput struct {
	// Error has a reference to an error if removing ingress rules fails.
	Error error
}

// RemoveIngressRules is a generic job used to remove ingress rules.
var RemoveIngressRules = &actions.Job{
	Execute: removeIngressRules,
}

// configureIngress is used by the ConfigureIngress job as the execute function.
func removeIngressRules(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	input := value.(RemoveIngressRulesInput)

	res, err := s.Platform().Orchestrator().Ingresses().Get(context.TODO(), input.Name, input.Namespace)
	if err != nil {
		return RemoveIngressRulesOutput{
			Error: err,
		}, nil
	}

	now := time.Now()
	timeout := s.Platform().Store().Orchestrator().Timeout()
	freq := s.Platform().Store().Orchestrator().PollFrequency()

	for t := now.Add(timeout); t.After(time.Now()); time.Sleep(freq) {
		// Exponential backoff
		freq *= 2

		var rule ingresses.Rule
		rule, err = s.Platform().Orchestrator().IngressRules().Get(context.TODO(), res, input.Host)
		if err != nil {
			continue
		}

		err = s.Platform().Orchestrator().IngressRules().Remove(context.TODO(), rule, input.Paths...)
		if err != nil {
			continue
		}

		break
	}

	return RemoveIngressRulesOutput{
		Error: err,
	}, nil
}
