package orchestrator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"io"
)

// Pods groups a set of methods to perform an operation with a Pod.
type Pods interface {
	Exec(resource Resource) Executor
	Reader(resource Resource) Reader
	WaitForCondition(resource Resource, condition Condition) waiter.Waiter
}

// Executor groups a set of methods to run commands and scripts inside a Pod.
type Executor interface {
	// Cmd runs a command inside a container.
	Cmd(command []string) error
	// Script runs a script inside a container.
	// Could be used to run copy_to_s3.sh
	Script(path string) error
}

// Reader groups a set of methods to read files and logs from a Pod.
type Reader interface {
	File(paths ...string) (io.Reader, error)
	Logs(container string, lines int64) (string, error)
}
