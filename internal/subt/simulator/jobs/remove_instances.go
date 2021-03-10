package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemoveInstances is a job in charge of removing all machines for a certain simulation.
var RemoveInstances = jobs.RemoveInstances.Extend(actions.Job{
	Name:       "remove-instances",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemoveInstancesInput},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
})

// prepareRemoveInstancesInput is in charge of preparing the input for the generic jobs.RemoveInstances job.
func prepareRemoveInstancesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	filters := make(map[string][]string)
	tags := subtapp.GetTagsInstanceBase(s.GroupID)

	for _, tag := range tags {
		for k, v := range tag.Map {
			filters[fmt.Sprintf("tag:%s", k)] = []string{v}
		}
	}

	return jobs.RemoveInstancesInput{
		machines.TerminateMachinesInput{
			Filters: filters,
		},
	}, nil
}
