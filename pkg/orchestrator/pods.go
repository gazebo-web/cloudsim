package orchestrator

import (
	"io"
)

// Pod groups a set of methods to identify a pod.
type Pod interface {
	Namespace() string
	Selector() string
	Name() string
}

// PodManager groups a set of methods to perform an operation with a Pod.
type PodManager interface {
	Exec(pod Resource) Executor
	Reader(pod Resource) Reader
	Condition(pod Resource, condition Condition) Waiter
}

// Executor groups a set of methods to run commands and scripts inside a Pod.
type Executor interface {
	// Cmd runs a command inside a container.
	Cmd(command []string)
	// Script runs a script inside a container.
	// Could be used to run copy_to_s3.sh
	Script(path string) error
}

// Reader groups a set of methods to read files and logs from a Pod.
type Reader interface {
	File(paths ...string) (io.Reader, error)
	Logs(lines int64) (string, error)
}
