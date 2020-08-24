package nodes

type Status string

const (
	StatusInitializing Status = "initializing"
	StatusRunning      Status = "running"
	StatusTerminating  Status = "terminating"
	StatusTerminated   Status = "terminated"
	StatusError        Status = "error"
)
