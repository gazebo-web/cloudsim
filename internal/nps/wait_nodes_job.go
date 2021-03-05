package nps

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForNodes is the job in charge of waiting until all the simulation instances have joined the cluster
// and have been marked as ready.
var WaitForNodes = jobs.Wait.Extend(actions.Job{
	Name:       "wait-for-nodes",
	PreHooks:   []actions.JobFunc{setStartState, createWaitForNodesInput},
	PostHooks:  []actions.JobFunc{checkWaitError, returnState},
	InputType:  actions.GetJobDataType(&StartSimulationData{}),
	OutputType: actions.GetJobDataType(&StartSimulationData{}),
})

// createWaitForNodesInput creates the input for the main wait job. This JobFunc is a pre-hook of the WaitForNodes job.
func createWaitForNodesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	startData := value.(*StartSimulationData)

	// Update the database entry with the latest status
	// \todo Help needed: I think this is not the recommended method to update
	// the database.
	var simEntry Simulation
	if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&simEntry).Error; err != nil {
		return nil, err
	}
	simEntry.Status = "Waiting for instance to join the K8 cluster."
	tx.Save(&simEntry)

	store.SetState(startData)

	// The NodeSelector was created in the LaunchInstances
	// \todo Improvment: Check that startData.NodeSelector was created.
	res := orchestrator.NewResource("", "", startData.NodeSelector)

	w := startData.Platform().Orchestrator().Nodes().WaitForCondition(res,
		orchestrator.ReadyCondition)

	// \todo Improvement: Each job has to return the correct type, but the
	// signature of the job functions call for an `interface{}` return type. This
	// make it very difficul to know that a function should return.
	return jobs.WaitInput{
		Request:       w,
		PollFrequency: startData.Platform().Store().Machines().PollFrequency(),
		Timeout:       startData.Platform().Store().Machines().Timeout(),
	}, nil
}
