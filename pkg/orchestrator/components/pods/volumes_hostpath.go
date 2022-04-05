package pods

import corev1 "k8s.io/api/core/v1"

// HostPathType defines the host path type used for volumes.
type HostPathType corev1.HostPathType

const (
	// HostPathUnset is used for backwards compatibility, leave it empty if unset.
	HostPathUnset = HostPathType(corev1.HostPathUnset)

	// HostPathDirectoryOrCreate should be set if nothing exists at the given path, an empty directory will be created
	// there as needed with file mode 0755.
	HostPathDirectoryOrCreate = HostPathType(corev1.HostPathDirectoryOrCreate)
)

// VolumeHostPath is a Volume implementation that mounts a directory in the host machine inside a container.
type VolumeHostPath struct {
	VolumeBase

	// HostPath represents a pre-existing file or directory on the host machine that is directly exposed to the
	// container.
	HostPath string

	// HostPathType defines the mount type and mounting behavior.
	HostPathType HostPathType
}
