package application

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"net/url"
	"strings"
)

// simPrefix is used to identify simulation pods.
const simPrefix = "sim"

// robotPrefix is used to identify robot simulation pods.
const robotPrefix = "rbt"

// GetPodNameFieldComputer is used to generate the name for a field computer pod for the given robot.
func GetPodNameFieldComputer(groupID simulations.GroupID, robotID string) string {
	return fmt.Sprintf("%s-%s-fc-%s", simPrefix, groupID, robotID)
}

// GetPodNameMoleBridge is used to generate the name for a mole bridge pod.
func GetPodNameMoleBridge(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-%s-mole-bridge", simPrefix, groupID)
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

// GetPodNameMappingServer is used to generate the name for the mapping server pod.
func GetPodNameMappingServer(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-%s-map-server", simPrefix, groupID)
}

// GetPodNameMappingServerCopy is used to generate the name for the mapping server copy pod.
func GetPodNameMappingServerCopy(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-copy", GetPodNameMappingServer(groupID))
}

// GetRobotID returns a robot identification name in the following form:
// rbtN with N being the given id.
// id requires that zero-indexes are used when calling GetRobotID.
func GetRobotID(id int) string {
	return fmt.Sprintf("%s%d", robotPrefix, id+1)
}

// GetContainerNameGazeboServer returns the gzserver container name.
func GetContainerNameGazeboServer() string {
	return "gzserver-container"
}

// GetContainerNameMoleBridge returns the Mole ROS/Pulsar bridge container name.
func GetContainerNameMoleBridge() string {
	return "mole-ros-pulsar-bridge"
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

// GetContainerNameMappingServer returns the mapping server container name.
func GetContainerNameMappingServer() string {
	return "mapping-server"
}

// GetContainerNameMappingServerCopy returns the mapping server copy container name.
func GetContainerNameMappingServerCopy() string {
	return "copy-to-s3"
}

// GetSimulationIngressPath gets the path to the gzserver websocket server for a certain simulation.
func GetSimulationIngressPath(groupID simulations.GroupID) string {
	return fmt.Sprintf("/simulations/%s", groupID.String())
}

// GetServiceNameWebsocket returns the websocket name for the given GroupID.
func GetServiceNameWebsocket(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-%s-websocket", simPrefix, groupID.String())
}

// GetCommsBridgeLogsFilename returns the filename for comms bridge logs.
func GetCommsBridgeLogsFilename(groupID simulations.GroupID, robotName string) string {
	return fmt.Sprintf("%s-fc-%s-commsbridge.tar.gz", groupID, strings.ToLower(robotName))
}

// GetGazeboLogsFilename returns the filename of the file that contains simulation logs.
func GetGazeboLogsFilename(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s.tar.gz", groupID.String())
}

// GetSimulationLogKey returns the path for logs inside a copy pod.
func GetSimulationLogKey(groupID simulations.GroupID, owner string) string {
	escaped := url.PathEscape(owner)
	return fmt.Sprintf("/gz-logs/%s/%s/", escaped, groupID)
}

func GetMappingServerLogKey(groupID simulations.GroupID, owner string) string {
	escaped := url.PathEscape(owner)
	return fmt.Sprintf("/map-logs/%s/%s/", escaped, groupID)
}
