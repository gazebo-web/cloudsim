package nps

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// WaitForPod waits for a pod to finish launching.
var WaitForPod = jobs.Wait.Extend(actions.Job{
	Name:       "wait-pod",
	PreHooks:   []actions.JobFunc{createWaitRequestForPod},
	PostHooks:  []actions.JobFunc{checkWaitError, returnState},
	InputType:  actions.GetJobDataType(&StartSimulationData{}),
	OutputType: actions.GetJobDataType(&StartSimulationData{}),
})

// createWaitRequestForPod is the pre hook in charge of passing the needed input to the Wait job.
func createWaitRequestForPod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	// Get the start simulation data for this job.
	startData := value.(*StartSimulationData)

	// Update the database entry with the latest status
	// \todo Help needed: I think this is not the recommended method to update
	// the database.
	var simEntry Simulation
	if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&simEntry).Error; err != nil {
		return nil, err
	}
	simEntry.Status = "Waiting for docker image (pod) IP."
	tx.Save(&simEntry)

	store.SetState(startData)

	// \todo What is the orchestrator? What is a Resource? Why do I have
	// create a new resource here? What is the name and namespace for the
	// resource (the documentation for orchestrator.NewResource doesn't say)?
	// What is a selector (the third parameter for NewResource)?
	//
	// I think the "selector" must match the "Labels" given to a Pod when it's
	// created?
	orchestratorResource := orchestrator.NewResource(
		"", "", startData.PodSelector)

	// Create wait for condition request
	req := startData.Platform().Orchestrator().Pods().WaitForCondition(
		orchestratorResource, orchestrator.HasIPStatusCondition)

	// Get timeout and poll frequency from store
	timeout := startData.Platform().Store().Machines().Timeout()
	pollFreq := startData.Platform().Store().Machines().PollFrequency()

	// Return new wait input
	return jobs.WaitInput{
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}
