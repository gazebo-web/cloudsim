package nps

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemovePods extends the generic jobs.RemovePods job. It's in charge of removing simulation pods.
var RemovePods = jobs.RemovePods.Extend(actions.Job{
	Name:       "remove-pods",
	PreHooks:   []actions.JobFunc{prepareRemovePodsInput},
	PostHooks:  []actions.JobFunc{checkRemovePodsNoError, returnState},
	InputType:  actions.GetJobDataType(&StopSimulationData{}),
	OutputType: actions.GetJobDataType(&StopSimulationData{}),
})

// checkRemovePodsNoError is a post-hook in charge of checking that no errors were thrown while removing pods.
// \todo Improvement: I blindly copy this from SubT and it seems to work. Move to a generic place so that apps don't have to keep copying the code.
func checkRemovePodsNoError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	stopData := store.State().(*StopSimulationData)
	out := value.(jobs.RemovePodsOutput)
	if out.Error != nil || len(out.Resources) != len(stopData.PodList) {
		err := deployment.SetJobData(tx, nil, actions.DeploymentJobData, out)
		if err != nil {
			return nil, err
		}
		return nil, out.Error
	}
	return nil, nil
}

// prepareRemovePodsInput is a pre-hook in charge of setting up the selector
// needed for the generic jobs to delete pods.
func prepareRemovePodsInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	fmt.Printf("\nRemoving a pod\n")

	stopData := store.State().(*StopSimulationData)

	// Update the database entry with the latest status
	// \todo Help needed: I think this is not the recommended method to update
	// the database.
	var sim Simulation
	if err := tx.Where("group_id = ?", stopData.GroupID.String()).First(&sim).Error; err != nil {
		return nil, err
	}
	sim.Status = "Removing docker image (pod)."
	tx.Save(&sim)

	// Namespace is the orchestrator namespace where simulations should be
	// launched.
	// \todo MAJOR ERROR: I would assume that this would return the value in
	// CLOUDSIM_MACHINES_ORCHESTRATOR_NAMESPACE. It is empty.
	namespace := stopData.Platform().Store().Orchestrator().Namespace()
	if namespace == "default" || namespace == "" {
		stopData.logger.Error("CLOUDSIM_ORCHESTRATOR_NAMESPACE has not been set")
		return nil, errors.New("CLOUDSIM_ORCHESTRATOR_NAMESPACE has not been set")
	}

	// Create a selector for the pod to remove
	labels := map[string]string{
		"cloudsim":         "true",
		"nps":              "true",
		"cloudsim_groupid": stopData.GroupID.String(),
	}
	stopData.PodSelector = orchestrator.NewSelector(labels)

	// Create a list of pods to remove. This application has only one.
	list := []orchestrator.Resource{
		orchestrator.NewResource(sim.Name, namespace, stopData.PodSelector),
	}

	// And if logs are enabled, gazebo server copy pod.
	// if s.Platform().Store().Ignition().LogsCopyEnabled() {
	// 	list = append(list, orchestrator.NewResource(subtapp.GetPodNameGazeboServerCopy(s.GroupID), ns, nil))
	// }

	stopData.PodList = list
	store.SetState(stopData)

	return jobs.RemovePodsInput(list), nil
}
