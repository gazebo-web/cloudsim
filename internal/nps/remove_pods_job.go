package nps

/*
import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemovePods extends the generic jobs.RemovePods job. It's in charge of removing simulation pods.
var RemovePods = jobs.RemovePods.Extend(actions.Job{
	Name:       "remove-pods",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemovePodsInput},
	PostHooks:  []actions.JobFunc{checkRemovePodsNoError, returnState},
	InputType:  actions.GetJobDataType(&StopSimulationData{}),
	OutputType: actions.GetJobDataType(&StopSimulationData{}),
})

// checkRemovePodsNoError is a post-hook in charge of checking that no errors were thrown while removing pods.
func checkRemovePodsNoError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)
	out := value.(jobs.RemovePodsOutput)
	if out.Error != nil || len(out.Resources) != len(s.PodList) {
		err := deployment.SetJobData(tx, nil, actions.DeploymentJobData, out)
		if err != nil {
			return nil, err
		}
		return nil, out.Error
	}
	return nil, nil
}

// prepareRemovePodsInput is a pre-hook in charge of setting up the selector needed for the generic jobs to delete pods.
func prepareRemovePodsInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	stopData := store.State().(*StopSimulationData)

  // \todo MAJOR ERROR: This returns empty
	// namespace := s.Platform().Store().Orchestrator().Namespace()
  namespace = "web-cloudsim-integration"

  var simEntry Simulation
  if err := tx.Where("group_id = ?", stopData.GroupID.String()).First(&simEntry).Error; err != nil {
    return nil, err
  }
  simEntry.Status = "Removing pod."
  tx.Save(&simEntry)


	// Remove the pod
	list := []orchestrator.Resource{
    orchestrator.NewResource(
      simEntry.Name,
      subtapp.GetPodNameGazeboServer(s.GroupID), namespace, nil)}

	stopData.PodList = list
	store.SetState(stopData)

	return jobs.RemovePodsInput(list), nil
}
*/
