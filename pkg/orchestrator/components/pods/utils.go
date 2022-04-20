package pods

import "fmt"

// NewChownContainer generates an init container that will be used to run the chown command in the /tmp folder.
func NewChownContainer(volumes []Volume, path string, uid, gid uint) Container {
	return Container{
		Name:    "chown-shared-volume",
		Image:   "infrastructureascode/aws-cli:latest",
		Command: []string{"/bin/sh"},
		Args:    []string{"-c", fmt.Sprintf("chown %d:%d %s", uid, gid, path)},
		Volumes: volumes,
	}
}
