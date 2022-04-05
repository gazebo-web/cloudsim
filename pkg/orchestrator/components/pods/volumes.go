package pods

// Volume is a storage that can be mounted to allow a Container to access its data.
type Volume interface {
	// Base contains generic information about the volume.
	Base() VolumeBase
}

// VolumeBase implements a base Volume.
// It is intended to be embedded in Volume type implementations.
type VolumeBase struct {
	// Name contains the name of the volume.
	Name string

	// MountPath is the path within the container at which the volume should be mounted.
	MountPath string

	// SubPath is an optional path within the mounted volume from which the container's volume should be mounted.
	// If not defined, the volume root will be used.
	SubPath string

	// ReadOnly indicates whether the volume should be mounted in read-only mode. Defaults to false.
	ReadOnly bool
}

// Base contains generic information about the volume.
func (v VolumeBase) Base() VolumeBase {
	return v
}
