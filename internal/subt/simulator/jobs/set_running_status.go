package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// SetSimulationStatusToRunning is used to set a simulation status to running.
var SetSimulationStatusToRunning = jobs.SetSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-wait-instances",
	PreHooks:   []actions.JobFunc{setRunningStatus},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

func setRunningStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	return jobs.SetSimulationStatusInput{
		GroupID: s.GroupID,
		Status:  simulations.StatusRunning,
	}, nil
}
