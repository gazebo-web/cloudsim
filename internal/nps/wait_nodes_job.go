package nps

import (
  "fmt"
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
  fmt.Printf("\n\ncreateWaitForNodesInput\n\n")
	startData := value.(*StartSimulationData)

  var simEntry Simulation
  if err := tx.Where("group_id = ?", startData.GroupID.String()).First(&simEntry).Error; err != nil {
    return nil, err
  }
  simEntry.Status = "Waiting for instance to join the K8 cluster."
  tx.Save(&simEntry)

	store.SetState(startData)

  // \todo: What is this and how do I choose a good value?
  selector := orchestrator.NewSelector(map[string]string{
    "cloudsim_groupid": startData.GroupID.String(),
  })

	res := orchestrator.NewResource("", "", selector)

	w := startData.Platform().Orchestrator().Nodes().WaitForCondition(res, orchestrator.ReadyCondition)

  fmt.Printf("\n\ndone createWaitForNodesInput\n\n")
  // \todo Sequencing jobs in the current format requires passing around interfaces of the correct type. If a type is not correct, the build is fine but the job sequence will error out. It's very difficult to figure out where the problem is.
  // It would be nice to have some kind of type of compile-time checking.
	return jobs.WaitInput{
		Request:       w,
		PollFrequency: startData.Platform().Store().Machines().PollFrequency(),
		Timeout:       startData.Platform().Store().Machines().Timeout(),
	}, nil
}
