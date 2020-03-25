package platform

type IPlatformContainers interface {
	GetGazeboServerName() string
	GetCommsBridgeName() string
	GetFieldComputerName() string
	GetSidecarName() string
}

func (p Platform) GetGazeboServerName() string {
	return "gzserver-container"
}

func (p Platform) GetCommsBridgeName() string {
	return "comms-bridge"
}

func (p Platform) GetFieldComputerName() string {
	return "field-computer"
}

func (p Platform) GetSidecarName() string {
	return "copy-to-s3"
}
