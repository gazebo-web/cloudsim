package jobs

import (
	"context"
	"github.com/gazebo-web/cloudsim/pkg/actions"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/ingresses"
	"github.com/gazebo-web/cloudsim/pkg/simulator/state"
	"github.com/jinzhu/gorm"
	"time"
)

// ConfigureIngressInput is the input for the ConfigureIngress job.
type ConfigureIngressInput struct {
	Name      string
	Namespace string
	Host      string
	Paths     []ingresses.Path
}

// ConfigureIngressOutput is the output of the ConfigureIngress job.
type ConfigureIngressOutput struct {
	// Error has a reference to an error if configuring the ingress fails.
	Error error
}

// ConfigureIngress is a generic job to be used to configure the ingress that will allow websocket connections.
var ConfigureIngress = &actions.Job{
	Execute: configureIngress,
}

// configureIngress is used by the ConfigureIngress job as the execute function.
func configureIngress(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	input := value.(ConfigureIngressInput)

	res, err := s.Platform().Orchestrator().Ingresses().Get(context.TODO(), input.Name, input.Namespace)
	if err != nil {
		return ConfigureIngressOutput{
			Error: err,
		}, nil
	}

	now := time.Now()
	timeout := s.Platform().Store().Orchestrator().Timeout()
	freq := s.Platform().Store().Orchestrator().PollFrequency()

	for t := now.Add(timeout); t.After(time.Now()); time.Sleep(freq) {
		// Exponential backoff
		freq *= 2

		// Get the rule for the given host
		var rule ingresses.Rule
		rule, err = s.Platform().Orchestrator().IngressRules().Get(context.TODO(), res, input.Host)
		if err != nil {
			continue
		}

		// Update paths.
		err = s.Platform().Orchestrator().IngressRules().Upsert(context.TODO(), rule, input.Paths...)
		if err != nil {
			continue
		}

		break
	}

	return ConfigureIngressOutput{
		Error: err,
	}, nil
}
