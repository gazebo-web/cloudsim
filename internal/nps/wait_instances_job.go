package nps

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForInstances is the job in charge of waiting for instances to be OK.
var WaitForInstances = jobs.WaitForInstances.Extend(actions.Job{
	Name:       "wait-for-instances",
	PreHooks:   []actions.JobFunc{createWaitForInstancesInput},
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&StartSimulationData{}),
	OutputType: actions.GetJobDataType(&StartSimulationData{}),
})

// createWaitForInstancesInput is the pre hook in charge of passing the list of created instances to the execute function.
func createWaitForInstancesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	startData := value.(*StartSimulationData)

  var simEntry Simulation
  if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&simEntry).Error; err != nil {
    return nil, err
  }
  simEntry.Status = "Waiting for cloud instances to launch."
  tx.Save(&simEntry)

	store.SetState(startData)
  // s := store.State().(*StartSimulationData)

	return jobs.WaitForInstancesInput(startData.CreateMachinesOutput), nil
}
