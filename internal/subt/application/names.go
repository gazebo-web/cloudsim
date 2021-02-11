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

// GetPodNameGazeboServerCopy is used to generate the name for the gzserver copy pod.
func GetPodNameGazeboServerCopy(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-copy", GetPodNameGazeboServer(groupID))
}

// GetPodNameGazeboServer is used to generate the name for the gazebo server pod.
func GetPodNameGazeboServer(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-%s-gzserver", simPrefix, groupID)
}

// GetRobotID returns a robot identification name in the following form:
// rbtN with N being the given id.
func GetRobotID(id int) string {
	return fmt.Sprintf("%s%d", robotPrefix, id+1)
}

// GetContainerNameGazeboServer returns the gzserver container name.
func GetContainerNameGazeboServer() string {
	return "gzserver-container"
}

// GetContainerNameCommsBridge returns the comms bridge container name.
func GetContainerNameCommsBridge() string {
	return "comms-bridge"
}

// GetContainerNameCommsBridgeCopy returns the comms bridge copy container name.
func GetContainerNameCommsBridgeCopy() string {
	return "copy-to-s3"
}

// GetContainerNameFieldComputer returns the field computer container name.
func GetContainerNameFieldComputer() string {
	return "field-computer"
}

// GetContainerNameGazeboServerCopy returns the gzserver copy container name.
func GetContainerNameGazeboServerCopy() string {
	return "copy-to-s3"
}

// GetSimulationIngressPath gets the path to the gzserver websocket server for a certain simulation.
func GetSimulationIngressPath(groupID simulations.GroupID) string {
	return fmt.Sprintf("/simulations/%s", groupID.String())
}

// GetServiceNameWebsocket returns the websocket name for the given GroupID.
func GetServiceNameWebsocket(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-%s-websocket", simPrefix, groupID)
}
