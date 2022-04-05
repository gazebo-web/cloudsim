package pods

// VolumeConfiguration is a Volume implementation that mounts the contents of a configuration inside a container.
type VolumeConfiguration struct {
	VolumeBase

	// ConfigurationName is the name of the configuration file to mount.
	ConfigurationName string

	// Items optionally defines the specific set of configuration entries to mount.
	// Keys indicate the items to mount, while values indicate where to mount them inside the volume.
	// If not defined the configuration will be mounted as a directory, with a file for each entry where the key of the
	// entry is the file name and the value the contents of the file.
	Items map[string]string
}
