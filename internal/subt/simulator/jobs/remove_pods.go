package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemovePods extends the generic jobs.RemovePods job. It's in charge of removing simulation pods.
var RemovePods = jobs.RemovePods.Extend(actions.Job{
	Name:       "remove-pods",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemovePodsInput},
	PostHooks:  []actions.JobFunc{checkRemovePodsNoError, returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
})

// checkRemovePodsNoError is a post-hook in charge of checking that no errors were thrown while removing pods.
func checkRemovePodsNoError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	out := value.(*jobs.RemovePodsOutput)
	if out.Error != nil {
		return nil, out.Error
	}
	return nil, nil
}

// prepareRemovePodsInput is a pre-hook in charge of setting up the selector needed for the generic jobs to delete pods.
func prepareRemovePodsInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	selector := subtapp.GetPodLabelsBase(s.GroupID, nil)

	return jobs.RemovePodsInput(selector), nil
}
