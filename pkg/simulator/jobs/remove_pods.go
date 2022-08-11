package jobs

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// RemovePodsInput is the input of the RemovePods job.
type RemovePodsInput []resource.Resource

// RemovePodsOutput is the output of the RemovePods job.
// This struct was set in place to let the post-hook handle errors.
type RemovePodsOutput struct {
	Resources []resource.Resource
	Error     error
}

// RemovePods is a generic job to remove pods from a cluster.
var RemovePods = &actions.Job{
	Execute: removePods,
}

// removePods is the main function executed by the RemovePods job.
func removePods(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	// Parse input
	input, ok := value.(RemovePodsInput)
	if !ok {
		return nil, simulator.ErrInvalidInput
	}

	if len(input) == 0 {
		return RemovePodsOutput{
			Resources: []resource.Resource{},
			Error:     nil,
		}, nil
	}

	var deleted []resource.Resource
	var err error

	for _, in := range input {
		var res resource.Resource
		res, err = s.Platform().Orchestrator().Pods().Delete(context.Background(), in)
		if err != nil {
			return nil, err
		}
		deleted = append(deleted, res)
	}

	return RemovePodsOutput{
		Resources: deleted,
		Error:     err,
	}, nil
}
