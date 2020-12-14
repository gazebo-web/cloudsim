package application

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

// GetEnvVarsCommsBridge returns the env vars for the comms-bridge container.
func GetEnvVarsCommsBridge(groupID simulations.GroupID, robotName, gzServerIP, verbosity string) map[string]string {
	return map[string]string{
		"IGN_PARTITION":  groupID.String(),
		"IGN_RELAY":      gzServerIP,
		"IGN_VERBOSE":    verbosity,
		"ROBOT_NAME":     robotName,
		"IGN_IP":         "", // To be removed.
		"ROS_MASTER_URI": "http://($ROS_IP):11311",
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

// GetEnvVarsCommsBridgeCopy returns the env vars for the comms-bridge-copy container.
func GetEnvVarsGazeboServerCopy(region, accessKey, secret string) map[string]string {
	return map[string]string{
		"AWS_DEFAULT_REGION":    region,
		"AWS_REGION":            region,
		"AWS_ACCESS_KEY_ID":     accessKey,
		"AWS_SECRET_ACCESS_KEY": secret,
	}
}
