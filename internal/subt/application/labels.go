package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"strings"
)

func GetNodeLabelsFieldComputer(groupID simulations.GroupID, robot simulations.Robot) orchestrator.Selector {
	return orchestrator.NewSelector(map[string]string{
		"cloudsim_groupid": groupID.String(),
		"field-computer":   "true",
		"robot_name":       strings.ToLower(robot.Name()),
	})
}

func GetNodeLabelsGazeboServer(groupID simulations.GroupID) orchestrator.Selector {
	return orchestrator.NewSelector(map[string]string{
		"cloudsim_groupid": groupID.String(),
		"gzserver":         "true",
	})
}

func GetPodLabelsFieldComputer(groupID simulations.GroupID, parent *simulations.GroupID) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		"field-computer": "true",
	})
	return base.Extend(ext)
}

func GetPodLabelsCommsBridge(groupID simulations.GroupID, parent *simulations.GroupID, robot simulations.Robot) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		"comms-bridge":    "true",
		"comms-for-robot": strings.ToLower(robot.Name()),
	})
	return base.Extend(ext)
}

func GetPodLabelsGazeboServer(groupID simulations.GroupID, parent *simulations.GroupID) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		"gzserver": "true",
	})
	return base.Extend(ext)
}

func getPodLabelsBase(groupID simulations.GroupID, parent *simulations.GroupID) orchestrator.Selector {
	base := orchestrator.NewSelector(map[string]string{
		"cloudsim":          "true",
		"SubT":              "true",
		"cloudsim-group-id": groupID.String(),
	})

	if parent != nil {
		base.Set("parent-group-id", parent.String())
	}

	return base
}
