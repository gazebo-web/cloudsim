package pods

import "fmt"

// NewChownContainer generates an init container that will be used to run the chown command in the /tmp folder.
func NewChownContainer(volumes []Volume) Container {
	return Container{
		Name:    "chown-shared-volume",
		Image:   "infrastructureascode/aws-cli:latest",
		Command: []string{"/bin/sh"},
		Args:    []string{"-c", fmt.Sprintf("chown %d:%d /tmp", 1000, 1000)},
		Volumes: volumes,
	}
}
