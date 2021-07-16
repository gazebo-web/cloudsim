package application

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// GetEnvVarsGazeboServer returns the env vars for the gazebo server container.
func GetEnvVarsGazeboServer(groupID simulations.GroupID, ip string, verbosity string) map[string]string {
	return map[string]string{
		"DISPLAY":          ":0",
		"QT_X11_NO_MITSHM": "1",
		"XAUTHORITY":       "/tmp/.docker.xauth",
		"USE_XVFB":         "1",
		"IGN_RELAY":        ip, // IP Cloudsim
		"IGN_PARTITION":    groupID.String(),
		"IGN_VERBOSE":      verbosity,
		"ROS_MASTER_URI":   "http://$(ROS_IP):11311",
	}
}

// GetEnvVarsFromSourceGazeboServer returns the env vars that need to be configured from a certain source for the gazebo
// server container.
func GetEnvVarsFromSourceGazeboServer() map[string]string {
	return map[string]string{
		"IGN_RELAY": pods.EnvVarSourcePodIP,
		"ROS_IP":    pods.EnvVarSourcePodIP,
		"IGN_IP":    pods.EnvVarSourcePodIP,
	}
}

// GetEnvVarsCommsBridge returns the env vars for the comms-bridge container.
func GetEnvVarsCommsBridge(groupID simulations.GroupID, robotName, gzServerIP, mappingServerIP, verbosity string) map[string]string {
	return map[string]string{
		"IGN_PARTITION":  groupID.String(),
		"IGN_RELAY":      fmt.Sprintf("%s:%s", gzServerIP, mappingServerIP),
		"IGN_VERBOSE":    verbosity,
		"ROBOT_NAME":     robotName,
		"ROS_MASTER_URI": "http://$(ROS_IP):11311",
	}
}

// GetEnvVarsFromSourceCommsBridge creates a map of the different env vars that should be configured from an external source.
// The resultant map will result in:
// "ENV_VAR_NAME": "SOURCE"
func GetEnvVarsFromSourceCommsBridge() map[string]string {
	return map[string]string{
		"ROS_IP": pods.EnvVarSourcePodIP,
		"IGN_IP": pods.EnvVarSourcePodIP,
	}
}

// GetEnvVarsFieldComputer returns the env vars for the field computer container.
func GetEnvVarsFieldComputer(robotName string, commsBridgeIP string) map[string]string {
	return map[string]string{
		"ROBOT_NAME":     robotName,
		"ROS_MASTER_URI": fmt.Sprintf("http://%s:11311", commsBridgeIP),
	}
}

// GetEnvVarsFromSourceFieldComputer returns the env vars for the field computer container.
func GetEnvVarsFromSourceFieldComputer() map[string]string {
	return map[string]string{
		"ROS_IP": pods.EnvVarSourcePodIP,
	}
}

// GetEnvVarsMappingServer returns the env vars for the mapping server container.
func GetEnvVarsMappingServer(groupID simulations.GroupID, gzServerIP string) map[string]string {
	return map[string]string{
		"IGN_PARTITION":  groupID.String(),
		"IGN_RELAY":      gzServerIP,
		"ROS_MASTER_URI": "http://$(ROS_IP):11311",
	}
}

// GetEnvVarsFromSourceMappingServer returns the env vars for the mapping server container from a certain source.
func GetEnvVarsFromSourceMappingServer() map[string]string {
	return map[string]string{
		"ROS_IP": pods.EnvVarSourcePodIP,
		"IGN_IP": pods.EnvVarSourcePodIP,
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

// GetEnvVarsMappingServerCopy returns the env vars for the mapping server copy container.
func GetEnvVarsMappingServerCopy(region, accessKey, secret string) map[string]string {
	return map[string]string{
		"AWS_DEFAULT_REGION":    region,
		"AWS_REGION":            region,
		"AWS_ACCESS_KEY_ID":     accessKey,
		"AWS_SECRET_ACCESS_KEY": secret,
	}
}
