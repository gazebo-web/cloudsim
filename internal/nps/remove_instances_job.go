package nps

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// RemoveInstances is a job in charge of removing all machines for a certain simulation.
var RemoveInstances = jobs.RemoveInstances.Extend(actions.Job{
	Name:       "remove-instances",
	PreHooks:   []actions.JobFunc{prepareRemoveInstancesInput},
	PostHooks:  []actions.JobFunc{saveState, returnState},
	InputType:  actions.GetJobDataType(&StopSimulationData{}),
	OutputType: actions.GetJobDataType(&StopSimulationData{}),
})

// prepareRemoveInstancesInput is in charge of preparing the input for the generic jobs.RemoveInstances job.
func prepareRemoveInstancesInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	stopData := store.State().(*StopSimulationData)

	// Update the database entry with the latest status
	// \todo Help needed: I think this is not the recommended method to update
	// the database.
	var sim Simulation
	if err := tx.Where("group_id = ?", stopData.GroupID.String()).First(&sim).Error; err != nil {
		return nil, err
	}
	sim.Status = "removing-instance"
	tx.Save(&sim)

	filters := make(map[string][]string)

	clusterName := stopData.Platform().Store().Machines().ClusterName()
	clusterKey := "kubernetes.io/cluster/" + clusterName

	// These are the tags to apply the EC2 machines.
	tags := []cloud.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"Name":                 sim.Name,
				"cloudsim_groupid":     string(stopData.GroupID),
				"project":              "nps",
				"Cloudsim":             "true",
				"cloudsim-application": "nps",

				// Note: `clusterKey` is extremely important. Without it, the EC2 node
				// will not join the cluster.
				clusterKey: "owned",
			},
		},
	}

	for _, tag := range tags {
		for k, v := range tag.Map {
			filters[fmt.Sprintf("tag:%s", k)] = []string{v}
		}
	}

	return jobs.RemoveInstancesInput{
		cloud.TerminateMachinesInput{
			Filters: filters,
		},
	}, nil
}

// saveState saves the simulation state.
func saveState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	stopData := store.State().(*StopSimulationData)

	// Update the database entry with the latest status
	// \todo Help needed: I think this is not the recommended method to update
	// the database.
	var sim Simulation
	if err := tx.Where("group_id = ?", stopData.GroupID.String()).First(&sim).Error; err != nil {
		return nil, err
	}
	sim.Status = "stopped"
	tx.Save(&sim)

	store.SetState(stopData)
	return stopData, nil
}
