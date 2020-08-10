package orchestrator

import (
	"bytes"
	"time"
)

// Condition represents a state that should be reached.
type Condition string

var (
	// ReadyCondition is a Condition set at ready.
	ReadyCondition Condition = "Ready"
)

// Orchestrator groups a set of methods to handle Nodes and Pods.
type Orchestrator interface {
	Nodes() NodeManager
	Pods() PodManager
}

// NodeManager groups a set of methods to register nodes into a cluster.
type NodeManager interface {
	Waiter
}

// PodManager groups a set of methods to perform operation with a pod.
type PodManager interface {
	Exec(selector string) Executor
	Reader(selector string) Reader
	Waiter
}

// Executor groups a set of methods to execute commands or scripts inside a pod.
type Executor interface {
	// Cmd runs a command inside a container.
	Cmd(command []string)
	// Script runs a script inside a container.
	// Could be used to run copy_to_s3.sh
	Script(path string) error
}

// Reader groups a set of methods to read a file or the logs from a pod.
type Reader interface {
	File(paths ...string) (*bytes.Buffer, error)
	Logs(lines int64) (*string, error)
}

// Waiter has a method to wait for a certain Condition to happen.
type Waiter interface {
	Wait(selector string, condition Condition, timeout time.Duration, pollFrequency time.Duration)
}
