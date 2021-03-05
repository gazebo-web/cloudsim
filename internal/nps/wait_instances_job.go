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
	// Get the start simulation data for this job.
	startData := value.(*StartSimulationData)

	// Update the database entry with the latest status
	// \todo Help needed: I think this is not the recommended method to update
	// the database.
	var simEntry Simulation
	if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&simEntry).Error; err != nil {
		return nil, err
	}
	simEntry.Status = "Waiting for cloud instances to launch."
	tx.Save(&simEntry)

	store.SetState(startData)

	return jobs.WaitForInstancesInput(startData.CreateMachinesOutput), nil
}
