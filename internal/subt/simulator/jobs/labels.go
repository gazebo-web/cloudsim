package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// JobPrepareLabels is in charge of preparing the labels for nodes and pods that will be used by the orchestrator.
var JobPrepareLabels = actions.Job{
	Name: "prepare-labels",
	PreHooks: []actions.JobFunc{
		setStartState,
		getStartGazeboNodeLabels,
		getStartFieldComputerNodeLabels,
		getStartGzServerPodLabels,
		getStartFieldComputerPodLabels,
		getStartBridgePodLabels,
		setParentGroupIDLabels,
	},
	Execute:         setStartState,
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: nil,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
}

func getStartGazeboNodeLabels(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	s.GazeboNodeLabels = map[string]string{
		"cloudsim_groupid": string(s.GroupID),
		"gzserver":         "true",
	}

	return s, nil
}

func getStartFieldComputerNodeLabels(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	s.FieldComputerNodeLabels = map[string]string{
		"cloudsim_groupid": string(s.GroupID),
		"field-computer":   "true",
	}

	return s, nil
}

func getStartGzServerPodLabels(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	s.GazeboServerPodLabels = map[string]string{
		"cloudsim":          "true",
		"SubT":              "true",
		"cloudsim-group-id": string(s.GroupID),
		"gzserver":          "true",
	}

	return s, nil
}

func getStartFieldComputerPodLabels(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	s.FieldComputerPodLabels = map[string]string{
		"cloudsim":          "true",
		"SubT":              "true",
		"cloudsim-group-id": string(s.GroupID),
		"field-computer":    "true",
	}

	return s, nil
}

func getStartBridgePodLabels(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	s.CommsBridgePodLabels = map[string]string{
		"cloudsim":          "true",
		"SubT":              "true",
		"cloudsim-group-id": string(s.GroupID),
		"comms-bridge":      "true",
	}

	return s, nil
}

func setParentGroupIDLabels(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	if sim.GetKind() == simulations.SimChild {
		parent, err := s.Services().Simulations().GetParent(s.GroupID)
		if err != nil {
			return nil, err
		}
		s.GazeboServerPodLabels["parent-group-id"] = string(parent.GetGroupID())
		s.FieldComputerPodLabels["parent-group-id"] = string(parent.GetGroupID())
		s.CommsBridgePodLabels["parent-group-id"] = string(parent.GetGroupID())
	}

	return s, nil
}
