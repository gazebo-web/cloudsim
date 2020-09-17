package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// SetSimulationStatusToLaunchInstances is used to set a simulation status to launch instances.
var SetSimulationStatusToLaunchInstances = jobs.SetSimulationStatus.Extend(actions.Job{
	Name:       "set-simulation-status-launch-instances",
	PreHooks:   []actions.JobFunc{setLaunchInstancesStatus},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

func setLaunchInstancesStatus(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	return jobs.SetSimulationStatusInput{
		GroupID: s.GroupID,
		Status:  simulations.StatusLaunchingInstances,
	}, nil
}
