package nps

import (
  "fmt"
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
  fmt.Printf("\n\n11\n\n")
	s := value.(*StartSimulationData)
  fmt.Printf("\n\n12\n\n")

	store.SetState(s)
  // s := store.State().(*StartSimulationData)
  fmt.Printf("\n\n13\n\n")

  fmt.Println("-------------------------------------------------------")
  fmt.Println(s.CreateMachinesInput)
  fmt.Println("-------------------------------------------------------")
	return jobs.WaitForInstancesInput(s.CreateMachinesOutput), nil
}
