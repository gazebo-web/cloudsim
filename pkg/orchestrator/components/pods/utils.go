package pods

import "fmt"

// NewChownContainer generates an init container that will be used to run the chown command in the /tmp folder.
func NewChownContainer(path string, uid, gid uint, volumes ...Volume) Container {
	return Container{
		Name:    "chown-shared-volume",
		Image:   "infrastructureascode/aws-cli:latest",
		Command: []string{"/bin/sh"},
		Args:    []string{"-c", fmt.Sprintf("chown -R %d:%d %s", uid, gid, path)},
		Volumes: volumes,
	}
}
