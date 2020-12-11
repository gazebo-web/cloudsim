package application

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// simPrefix is used to identify simulation pods.
const simPrefix = "sim"

// robotPrefix is used to identify robot simulation pods.
const robotPrefix = "rbt"

// GetPodNameFieldComputer is used to generate the name for a field computer pod for the given robot.
func GetPodNameFieldComputer(groupID simulations.GroupID, robotID string) string {
	return fmt.Sprintf("%s-%s-fc-%s", simPrefix, groupID, robotID)
}

// GetPodNameCommsBridge is used to generate the name for a comms bridge pod for the given robot.
func GetPodNameCommsBridge(groupID simulations.GroupID, robotID string) string {
	return fmt.Sprintf("%s-%s-comms-%s", simPrefix, groupID, robotID)
}

// GetPodNameCommsBridgeCopy is used to generate the name for the comms bridge copy pod.
func GetPodNameCommsBridgeCopy(groupID simulations.GroupID, robotID string) string {
	return fmt.Sprintf("%s-copy", GetPodNameCommsBridge(groupID, robotID))
}

// GetPodNameGazeboServer is used to generate the name for the gazebo server pod.
func GetPodNameGazeboServer(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-%s-gzserver", simPrefix, groupID)
}

// GetRobotID returns a robot identification name in the following form:
// rbtN with N being the given id.
func GetRobotID(id int) string {
	return fmt.Sprintf("%s%d", robotPrefix, id)
}

// GetContainerNameCommsBridge returns the comms bridge container name.
func GetContainerNameCommsBridge() string {
	return "comms-bridge"
}

// GetContainerNameCommsBridgeCopy returns the comms bridge copy container name.
func GetContainerNameCommsBridgeCopy() string {
	return "copy-to-s3"
}
