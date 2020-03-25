package platform

// IPlatformContainer represents a set of methods for containers.
type IPlatformContainer interface {
	GetGazeboServerName() string
	GetCommsBridgeName() string
	GetFieldComputerName() string
	GetSidecarName() string
}

// GetGazeboServerName returns the name for the gazebo server container.
func (p Platform) GetGazeboServerName() string {
	return "gzserver-container"
}

// GetCommsBridgeName returns the name for the comms bridge container.
func (p Platform) GetCommsBridgeName() string {
	return "comms-bridge"
}

// GetFieldComputerName returns the name for the field computer container.
func (p Platform) GetFieldComputerName() string {
	return "field-computer"
}

// GetSidecarName returns the name for the sidecar container.
func (p Platform) GetSidecarName() string {
	return "copy-to-s3"
}
