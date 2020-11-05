package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForInstances is the job in charge of waiting for instances to be OK.
var WaitForInstances = jobs.WaitForInstances.Extend(actions.Job{
	Name:       "wait-for-instances",
	PreHooks:   []actions.JobFunc{createWaitForInstancesInput},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createWaitForInstancesInput is the pre hook in charge of passing the list of created instances to the execute function.
func createWaitForInstancesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	return jobs.WaitForInstancesInput(s.CreateMachinesOutput), nil
}
