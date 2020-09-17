package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// SetSimulationStatusToWaitPods is used to set a simulation status to wait pods.
var SetSimulationStatusToWaitPods = jobs.SetSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-wait-pods",
	PreHooks:   []actions.JobFunc{setWaitPodsStatus},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

func setWaitPodsStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	return jobs.SetSimulationStatusInput{
		GroupID: s.GroupID,
		Status:  simulations.StatusWaitingPods,
	}, nil
}
