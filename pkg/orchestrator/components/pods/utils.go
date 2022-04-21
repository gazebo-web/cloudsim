package pods

import (
	"fmt"
	"strings"
)

// NewChownContainer generates an init container that sets the owner of volume mount directories to a specific user
// and group.
func NewChownContainer(uid, gid uint, volumes ...Volume) Container {
	mountPaths := make([]string, len(volumes))
	for i, volume := range volumes {
		mountPaths[i] = volume.Base().MountPath
	}
	paths := strings.Join(mountPaths, " ")

	return Container{
		Name:    "chown-shared-volume",
		Image:   "infrastructureascode/aws-cli:latest",
		Command: []string{"/bin/sh"},
		Args:    []string{"-c", fmt.Sprintf("chown -R %d:%d %s", uid, gid, paths)},
		Volumes: volumes,
	}
}
