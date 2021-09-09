package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
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
	labelMoleBridge          = "mole-bridge"
	labelCommsBridge         = "comms-bridge"
	labelMappingServer       = "mapping-server"
	labelCommsBridgeForRobot = "comms-for-robot"
	labelCopyS3              = "copy-to-s3"
	labelCopyForRobot        = "copy-for-robot"
	labelCloudsim            = "cloudsim"
	labelSubT                = "SubT"
)

// GetNodeLabelsFieldComputer returns a selector that identifies a field computer node.
func GetNodeLabelsFieldComputer(groupID simulations.GroupID, robot simulations.Robot) resource.Selector {
	base := GetNodeLabelsBase(groupID)

	return base.Extend(resource.NewSelector(map[string]string{
		labelFieldComputer: "true",
		labelRobotName:     strings.ToLower(robot.GetName()),
	}))
}

// GetNodeLabelsGazeboServer returns a selector that identifies a gazebo node.
func GetNodeLabelsGazeboServer(groupID simulations.GroupID) resource.Selector {
	base := GetNodeLabelsBase(groupID)

	return base.Extend(resource.NewSelector(map[string]string{
		labelGazeboServer: "true",
	}))
}

// GetNodeLabelsMappingServer returns a selector that identifies a mapping server node.
func GetNodeLabelsMappingServer(groupID simulations.GroupID) resource.Selector {
	base := GetNodeLabelsBase(groupID)

	return base.Extend(resource.NewSelector(map[string]string{
		labelMappingServer: "true",
	}))
}

// GetNodeLabelsBase returns the base labels to identify a simulation's node.
func GetNodeLabelsBase(groupID simulations.GroupID) resource.Selector {
	return resource.NewSelector(map[string]string{
		labelGroupID: groupID.String(),
	})
}

// GetPodLabelsFieldComputer returns a selector that identifies a field computer pod.
func GetPodLabelsFieldComputer(groupID simulations.GroupID, parent *simulations.GroupID) resource.Selector {
	base := GetPodLabelsBase(groupID, parent)
	ext := resource.NewSelector(map[string]string{
		labelFieldComputer: "true",
	})
	return base.Extend(ext)
}

// GetPodLabelsMoleBridge returns a selector that identifies a mole bridge pod.
func GetPodLabelsMoleBridge(groupID simulations.GroupID, parent *simulations.GroupID) resource.Selector {
	base := GetPodLabelsBase(groupID, parent)
	ext := resource.NewSelector(map[string]string{
		labelMoleBridge: "true",
	})
	return base.Extend(ext)
}

// GetPodLabelsCommsBridge returns a selector that identifies a comms bridge pod.
func GetPodLabelsCommsBridge(groupID simulations.GroupID, parent *simulations.GroupID, robot simulations.Robot) resource.Selector {
	base := GetPodLabelsBase(groupID, parent)
	ext := resource.NewSelector(map[string]string{
		labelCommsBridge:         "true",
		labelCommsBridgeForRobot: strings.ToLower(robot.GetName()),
	})
	return base.Extend(ext)
}

// GetPodLabelsCommsBridgeCopy returns a selector that identifies a comms bridge copy pod.
func GetPodLabelsCommsBridgeCopy(groupID simulations.GroupID, parent *simulations.GroupID, robot simulations.Robot) resource.Selector {
	base := GetPodLabelsBase(groupID, parent)
	ext := resource.NewSelector(map[string]string{
		labelCopyS3:       "true",
		labelCopyForRobot: strings.ToLower(robot.GetName()),
	})
	return base.Extend(ext)
}

// GetPodLabelsGazeboServerCopy returns a selector that identifies a gzserver copy pod.
func GetPodLabelsGazeboServerCopy(groupID simulations.GroupID, parent *simulations.GroupID) resource.Selector {
	base := GetPodLabelsBase(groupID, parent)
	ext := resource.NewSelector(map[string]string{
		labelCopyS3: "true",
	})
	return base.Extend(ext)
}

// GetPodLabelsGazeboServer returns a selector that identifies a gzserver pod.
func GetPodLabelsGazeboServer(groupID simulations.GroupID, parent *simulations.GroupID) resource.Selector {
	base := GetPodLabelsBase(groupID, parent)
	ext := resource.NewSelector(map[string]string{
		labelGazeboServer: "true",
	})
	return base.Extend(ext)
}

// GetPodLabelsBase returns the base set of key-values for all pod selectors.
func GetPodLabelsBase(groupID simulations.GroupID, parent *simulations.GroupID) resource.Selector {
	base := resource.NewSelector(map[string]string{
		labelCloudsim:   "true",
		labelSubT:       "true",
		labelPodGroupID: groupID.String(),
	})

	if parent != nil {
		base.Set(labelParentGroupID, parent.String())
	}

	return base
}

// GetWebsocketServiceLabels returns a selector that will identify a websocket service for a certain simulation.
func GetWebsocketServiceLabels(groupID simulations.GroupID) resource.Selector {
	return resource.NewSelector(map[string]string{
		labelPodGroupID: groupID.String(),
	})
}

// GetPodLabelsMappingServer returns a selector that identifies a mapping server pod pod.
func GetPodLabelsMappingServer(groupID simulations.GroupID, parent *simulations.GroupID) resource.Selector {
	base := GetPodLabelsBase(groupID, parent)
	ext := resource.NewSelector(map[string]string{
		labelMappingServer: "true",
	})
	return base.Extend(ext)
}

// GetPodLabelsMappingServerCopy returns a selector that identifies a mapping server copy pod.
func GetPodLabelsMappingServerCopy(groupID simulations.GroupID, parent *simulations.GroupID) resource.Selector {
	base := GetPodLabelsBase(groupID, parent)
	ext := resource.NewSelector(map[string]string{
		labelCopyS3: "true",
	})
	return base.Extend(ext)
}
