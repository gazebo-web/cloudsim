package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/context"
)

// PrepareLabels is a job in charge of preparing pod labels
var PrepareLabels = &actions.Job{
	Name:       "prepare-labels",
	Execute:    prepareLabels,
	InputType:  actions.GetJobDataType(simulations.GroupID("")),
	OutputType: actions.GetJobDataType(simulations.GroupID("")),
}

func prepareLabels(ctx actions.Context, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {

	data := ctx.Store().Get().(*StartSimulationData)

	data.GazeboNodeSelector = map[string]string{
		"cloudsim_groupid": string(data.GroupID),
		"gzserver":         "true",
	}

	// Set up pod labels
	data.GazeboLabels = map[string]string{
		"cloudsim":          "true",
		"SubT":              "true",
		"cloudsim-group-id": string(data.GroupID),
		"gzserver":          "true",
	}

	data.FieldComputerLabels = map[string]string{
		"cloudsim":          "true",
		"SubT":              "true",
		"cloudsim-group-id": string(data.GroupID),
		"field-computer":    "true",
	}

	data.BridgeLabels = map[string]string{
		"cloudsim":          "true",
		"SubT":              "true",
		"cloudsim-group-id": string(data.GroupID),
		"comms-bridge":      "true",
	}

	simCtx := context.NewContext(ctx)

	sim, err := simCtx.Services().Simulations().Get(data.GroupID)
	if err != nil {
		return nil, err
	}

	// If simulation is child, add another label with the parent's group id.
	if sim.Kind() == simulations.SimChild {
		parent, err := simCtx.Services().Simulations().GetParent(data.GroupID)
		if err != nil {
			return nil, err
		}
		data.GazeboLabels["parent-group-id"] = string(parent.GroupID())
		data.FieldComputerLabels["parent-group-id"] = string(parent.GroupID())
		data.BridgeLabels["parent-group-id"] = string(parent.GroupID())
	}

	err = ctx.Store().Set(data)
	if err != nil {
		return nil, err
	}

	return data.GroupID, nil
}
