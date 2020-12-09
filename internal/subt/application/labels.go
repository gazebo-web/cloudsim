package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"strings"
)

const (
	labelGroupID             = "cloudsim_groupid"
	labelPodGroupID          = "cloudsim-group-id"
	labelParentGroupID       = "parent-group-id"
	labelFieldComputer       = "field-computer"
	labelRobotName           = "robot_name"
	labelGazeboServer        = "gzserver"
	labelCommsBridge         = "comms-bridge"
	labelCommsBridgeForRobot = "comms-for-robot"
	labelCopyS3              = "copy-to-s3"
	labelCopyForRobot        = "copy-for-robot"
	labelCloudsim            = "cloudsim"
	labelSubT                = "SubT"
)

// GetNodeLabelsFieldComputer returns a selector that identifies a field computer node.
func GetNodeLabelsFieldComputer(groupID simulations.GroupID, robot simulations.Robot) orchestrator.Selector {
	return orchestrator.NewSelector(map[string]string{
		labelGroupID:       groupID.String(),
		labelFieldComputer: "true",
		labelRobotName:     strings.ToLower(robot.Name()),
	})
}

// GetNodeLabelsGazeboServer returns a selector that identifies a gazebo node.
func GetNodeLabelsGazeboServer(groupID simulations.GroupID) orchestrator.Selector {
	return orchestrator.NewSelector(map[string]string{
		labelGroupID:      groupID.String(),
		labelGazeboServer: "true",
	})
}

// GetPodLabelsFieldComputer returns a selector that identifies a field computer pod.
func GetPodLabelsFieldComputer(groupID simulations.GroupID, parent *simulations.GroupID) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		labelFieldComputer: "true",
	})
	return base.Extend(ext)
}

// GetPodLabelsCommsBridge returns a selector that identifies a comms bridge pod.
func GetPodLabelsCommsBridge(groupID simulations.GroupID, parent *simulations.GroupID, robot simulations.Robot) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		labelCommsBridge:         "true",
		labelCommsBridgeForRobot: strings.ToLower(robot.Name()),
	})
	return base.Extend(ext)
}

// GetPodLabelsCommsBridgeCopy returns a selector that identifies a comms bridge copy pod.
func GetPodLabelsCommsBridgeCopy(groupID simulations.GroupID, parent *simulations.GroupID, robot simulations.Robot) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		labelCopyS3:       "true",
		labelCopyForRobot: strings.ToLower(robot.Name()),
	})
	return base.Extend(ext)
}

// GetPodLabelsGazeboServer returns a selector that identifies a gzserver pod.
func GetPodLabelsGazeboServer(groupID simulations.GroupID, parent *simulations.GroupID) orchestrator.Selector {
	base := getPodLabelsBase(groupID, parent)
	ext := orchestrator.NewSelector(map[string]string{
		labelGazeboServer: "true",
	})
	return base.Extend(ext)
}

// getPodLabelsBase returns the base set of key-values for all pod selectors.
func getPodLabelsBase(groupID simulations.GroupID, parent *simulations.GroupID) orchestrator.Selector {
	base := orchestrator.NewSelector(map[string]string{
		labelCloudsim:   "true",
		labelSubT:       "true",
		labelPodGroupID: groupID.String(),
	})

	if parent != nil {
		base.Set(labelParentGroupID, parent.String())
	}

	return base
}