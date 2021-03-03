package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// RemoveIngressRulesInput is the input for the RemoveIngressRule job.
type RemoveIngressRulesInput struct {
	Name      string
	Namespace string
	Host      string
	Paths     []orchestrator.Path
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

	res, err := s.Platform().Orchestrator().Ingresses().Get(input.Name, input.Namespace)
	if err != nil {
		return RemoveIngressRulesOutput{Error: err}, nil
	}

	rule, err := s.Platform().Orchestrator().IngressRules().Get(res, input.Host)
	if err != nil {
		return RemoveIngressRulesOutput{Error: err}, nil
	}

	err = s.Platform().Orchestrator().IngressRules().Remove(rule, input.Paths...)
	if err != nil {
		return RemoveIngressRulesOutput{Error: err}, nil
	}

	return RemoveIngressRulesOutput{Error: nil}, nil
}
