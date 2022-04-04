package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// RemoveConfigurationsInput is the input for the RemoveConfigurations job.
type RemoveConfigurationsInput struct {
	Resource resource.Resource
}

// RemoveConfigurationsOutput is the output of the RemoveConfigurations job.
type RemoveConfigurationsOutput struct {
	// Error has a reference to the latest error thrown when removing the configurations.
	Error error
}

// RemoveConfigurations is a generic job that removes configurations.
var RemoveConfigurations = &actions.Job{
	Execute: removeConfigurations,
}

// removeConfigurations is used by the RemoveConfigurations job as the execute function.
func removeConfigurations(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	input := value.(RemoveConfigurationsInput)

	_, err := s.Platform().Orchestrator().Configurations().Delete(input.Resource)

	return RemoveConfigurationsOutput{
		Error: err,
	}, nil
}
