package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// SetSimulationStatusToWaitInstances is used to set a simulation status to wait instances.
var SetSimulationStatusToWaitInstances = jobs.SetSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-wait-instances",
	PreHooks:   []actions.JobFunc{setWaitInstancesStatus},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

func setWaitInstancesStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	return jobs.SetSimulationStatusInput{
		GroupID: s.GroupID,
		Status:  simulations.StatusWaitingInstances,
	}, nil
}
