package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemoveInstances is a job in charge of removing all machines for a certain simulation.
var RemoveInstances = jobs.RemoveInstances.Extend(actions.Job{
	Name:       "remove-instances",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemoveInstances},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
})

// prepareRemoveInstances is in charge of preparing the input for the generic jobs.RemoveInstances job.
func prepareRemoveInstances(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	return jobs.RemoveInstancesInput{
		cloud.TerminateMachinesInput{
			Filters: map[string][]string{
				fmt.Sprintf("tag:%s", "CloudsimGroupID"): {
					s.GroupID.String(),
				},
			},
		},
	}, nil
}
