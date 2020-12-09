package application

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// prefix is used to identify as a prefix for pod names to identify a simulation.
const prefix = "sim"

// GetFieldComputerPodName is used to generate the name for a field computer pod for the given robot.
func GetFieldComputerPodName(groupID simulations.GroupID, robotID string) string {
	return fmt.Sprintf("%s-%s-fc-%s", prefix, groupID, robotID)
}

// GetCommsBridgePodName is used to generate the name for a comms bridge pod for the given robot.
func GetCommsBridgePodName(groupID simulations.GroupID, robotID string) string {
	return fmt.Sprintf("%s-%s-comms-%s", prefix, groupID, robotID)
}

// GetGazeboServerPodName is used to generate the name for the gazebo server pod.
func GetGazeboServerPodName(groupID simulations.GroupID) string {
	return fmt.Sprintf("%s-%s-gzserver", prefix, groupID)
}

// GetSimulationIngressPath gets the path to the gzserver websocket server for a certain simulation.
func GetSimulationIngressPath(groupID simulations.GroupID) string {
	return fmt.Sprintf("/simulations/%s", groupID.String())
}
