package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"strings"
)

// GetNodeLabelsFieldComputer returns a selector that identifies a field computer node.
func GetNodeLabelsFieldComputer(groupID simulations.GroupID, robot simulations.Robot) orchestrator.Selector {
	return orchestrator.NewSelector(map[string]string{
		"cloudsim_groupid": groupID.String(),
		"field-computer":   "true",
		"robot_name":       strings.ToLower(robot.Name()),
	})
}

// GetNodeLabelsGazeboServer returns a selector that identifies a gazebo node.
func GetNodeLabelsGazeboServer(groupID simulations.GroupID) orchestrator.Selector {
	return orchestrator.NewSelector(map[string]string{
		"cloudsim_groupid": groupID.String(),
		"gzserver":         "true",
	})
}

// GetPodLabelsFieldComputer returns a selector that identifies a field computer pod.
func GetPodLabelsFieldComputer(groupID simulations.GroupID, parent *simulations.GroupID) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		"field-computer": "true",
	})
	return base.Extend(ext)
}

// GetPodLabelsCommsBridge returns a selector that identifies a comms bridge pod.
func GetPodLabelsCommsBridge(groupID simulations.GroupID, parent *simulations.GroupID, robot simulations.Robot) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		"comms-bridge":    "true",
		"comms-for-robot": strings.ToLower(robot.Name()),
	})
	return base.Extend(ext)
}

// GetPodLabelsCommsBridgeCopy returns a selector that identifies a comms bridge copy pod.
func GetPodLabelsCommsBridgeCopy(groupID simulations.GroupID, parent *simulations.GroupID, robot simulations.Robot) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		"copy-to-s3":     "true",
		"copy-for-robot": strings.ToLower(robot.Name()),
	})
	return base.Extend(ext)
}

// GetPodLabelsGazeboServer returns a selector that identifies a gzserver pod.
func GetPodLabelsGazeboServer(groupID simulations.GroupID, parent *simulations.GroupID) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		"gzserver": "true",
	})
	return base.Extend(ext)
}

// getPodLabelsBase returns the base set of key-values for all pod selectors.
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
