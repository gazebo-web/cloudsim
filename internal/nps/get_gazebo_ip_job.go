package nps

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// GetPodIP is a job in charge of getting the IP from a server pod.
var GetPodIP = &actions.Job{
	Name:       "get-pod-ip",
	PreHooks:   []actions.JobFunc{setStartState},
	Execute:    getIP,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&StartSimulationData{}),
	OutputType: actions.GetJobDataType(&StartSimulationData{}),
}

// getGazeboIP gets the gazebo server pod IP and assigns it to the start simulation state.
func getIP(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	startData := store.State().(*StartSimulationData)

	var simEntry Simulation
	if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&simEntry).Error; err != nil {
		return nil, err
	}

	// Namespace is the orchestrator namespace where simulations should be
	// launched.
	// \todo MAJOR ERROR: I would assume that this would return the value in
	// CLOUDSIM_MACHINES_ORCHESTRATOR_NAMESPACE. It is empty.
	namespace := startData.Platform().Store().Orchestrator().Namespace()
	if namespace == "default" || namespace == "" {
		startData.logger.Error("CLOUDSIM_ORCHESTRATOR_NAMESPACE has not been set")
		return nil, errors.New("CLOUDSIM_ORCHESTRATOR_NAMESPACE has not been set")
	}

	ip, err := startData.Platform().Orchestrator().Pods().GetIP(
		// This is the name of the Pod, set in launch_pod_job.go
		simEntry.Name,

		// This is set by the CLOUDSIM_ORCHESTRATOR_NAMESPACE
		namespace)

	if err != nil {
		return nil, err
	}

	simEntry.Status = "IP has been acquired."
	simEntry.IP = ip
	tx.Save(&simEntry)

	startData.IP = ip

	return startData, nil
}
