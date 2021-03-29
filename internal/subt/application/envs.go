package application

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// GetEnvVarsCommsBridge returns the env vars for the comms-bridge container.
func GetEnvVarsCommsBridge(groupID simulations.GroupID, robotName, gzServerIP, verbosity string) map[string]string {
	return map[string]string{
		"IGN_PARTITION":  groupID.String(),
		"IGN_RELAY":      gzServerIP,
		"IGN_VERBOSE":    verbosity,
		"ROBOT_NAME":     robotName,
		"IGN_IP":         "", // To be removed.
		"ROS_MASTER_URI": "http://$(ROS_IP):11311",
	}
}

// GetEnvVarsFromSourceCommsBridge creates a map of the different env vars that should be configured from an external source.
// The resultant map will result in:
// "ENV_VAR_NAME": "SOURCE"
func GetEnvVarsFromSourceCommsBridge() map[string]string {
	return map[string]string{
		"ROS_IP": orchestrator.EnvVarSourcePodIP,
	}
}

// GetEnvVarsFieldComputer returns the env vars for the field computer container.
func GetEnvVarsFieldComputer(robotName string, commsBridgeIP string) map[string]string {
	return map[string]string{
		"ROBOT_NAME":     robotName,
		"ROS_MASTER_URI": fmt.Sprintf("http://%s:11311", commsBridgeIP),
	}
}

// GetEnvVarsCommsBridgeCopy returns the env vars for the comms-bridge-copy container.
func GetEnvVarsCommsBridgeCopy(region, accessKey, secret string) map[string]string {
	return map[string]string{
		"AWS_DEFAULT_REGION":    region,
		"AWS_REGION":            region,
		"AWS_ACCESS_KEY_ID":     accessKey,
		"AWS_SECRET_ACCESS_KEY": secret,
	}
}

// GetEnvVarsGazeboServerCopy returns the env vars for the gzserver copy container.
func GetEnvVarsGazeboServerCopy(region, accessKey, secret string) map[string]string {
	return map[string]string{
		"AWS_DEFAULT_REGION":    region,
		"AWS_REGION":            region,
		"AWS_ACCESS_KEY_ID":     accessKey,
		"AWS_SECRET_ACCESS_KEY": secret,
	}
}
