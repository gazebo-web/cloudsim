package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// getDataFromJob is used to get the current job data from context.
func getDataFromJob(ctx actions.Context, deployment *actions.Deployment) (interface{}, error) {
	simCtx := context.NewContext(ctx)

	data := simCtx.Value(deployment.CurrentJob).(*StartSimulationData)

	return data, nil
}
