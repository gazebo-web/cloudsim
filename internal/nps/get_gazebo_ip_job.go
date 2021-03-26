package nps

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
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

	// Get the address of the node, which is publically accessible.
	address, err := startData.Platform().Orchestrator().Nodes().GetExternalDNSAddress(orchestrator.NewResource("", "", startData.NodeSelector))

	if err != nil {
		return nil, err
	}
	if len(address) == 0 || address == nil {
		return nil, errors.New("Unable to get node address")
	}

	simEntry.Status = "Address has been acquired."
	simEntry.URI = "http://" + address[0] + ":8080/vnc.html"
	tx.Save(&simEntry)

	startData.URI = simEntry.URI

	return startData, nil
}
