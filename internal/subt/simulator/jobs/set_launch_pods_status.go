package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// SetSimulationStatusToLaunchPods is used to set a simulation status to launch pods.
var SetSimulationStatusToLaunchPods = jobs.SetSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-launch-pods",
	PreHooks:   []actions.JobFunc{setLaunchPodsStatus},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

func setLaunchPodsStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	return jobs.SetSimulationStatusInput{
		GroupID: s.GroupID,
		Status:  simulations.StatusLaunchingPods,
	}, nil
}
