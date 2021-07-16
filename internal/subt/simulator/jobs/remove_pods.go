package jobs

import (
	"github.com/jinzhu/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	simulationspkg "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemovePods extends the generic jobs.RemovePods job. It's in charge of removing simulation pods.
var RemovePods = jobs.RemovePods.Extend(actions.Job{
	Name:       "remove-pods",
	PreHooks:   []actions.JobFunc{setStopState, prepareRemovePodsInput},
	PostHooks:  []actions.JobFunc{checkRemovePodsNoError, returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
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
	s := store.State().(*state.StopSimulation)

	sim, err := s.SubTServices().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	var parentGroupID *simulationspkg.GroupID
	if sim.IsKind(simulationspkg.SimParent) {
		parentSim, err := s.SubTServices().Simulations().GetParent(s.GroupID)
		if err != nil {
			return nil, err
		}
		groupID := parentSim.GetGroupID()
		parentGroupID = &groupID
	}

	ns := s.Platform().Store().Orchestrator().Namespace()

	pods, err := s.Platform().Orchestrator().Pods().List(ns, subtapp.GetPodLabelsBase(sim.GetGroupID(), parentGroupID))
	if err != nil {
		return nil, err
	}

	podList := make([]resource.Resource, len(pods))
	for i, pod := range pods {
		podList[i] = pod
	}

	s.PodList = podList
	store.SetState(s)

	return jobs.RemovePodsInput(podList), nil
}
