package simulations

type Status int

const (
	StatusPending        Status = 0
	StatusLaunchingNodes Status = 10
	StatusLaunchingPods  Status = 20
	// StatusParentLaunching is only used for Parent simulations when some of their children
	// are still launching and there wasn't any errors so far.
	StatusParentLaunching Status = 25
	// StatusParentLaunchingWithErrors is only used for Parent simulations when some of their children
	// finished with errors and some are still launching/running.
	StatusParentLaunchingWithErrors Status = 28
	StatusRunning                   Status = 30
	// StatusRunningWithErrors is used for Parent simulations when some of their children
	// finished with errors and some are still running.
	// @deprecated do not use.
	StatusRunningWithErrors    Status = 40
	StatusTerminateRequested   Status = 50
	StatusDeletingPods         Status = 60
	StatusDeletingNodes        Status = 70
	StatusTerminatingInstances Status = 80
	StatusTerminated           Status = 90
	StatusRejected             Status = 100
)

// GetStatusLabel returns a status string from the given status.
func GetStatusLabel(status Status) string {
	switch status {
	case StatusPending:
		return "Pending"
		break
	case StatusLaunchingNodes:
		return "LaunchingNodes"
		break
	case StatusLaunchingPods:
		return "LaunchingPods"
		break
	case StatusParentLaunching:
		return "Launching"
		break
	case StatusParentLaunchingWithErrors:
		return "RunningWithErrors"
		break
	case StatusRunning:
		return "Running"
		break
	case StatusRunningWithErrors:
		return "RunningWithErrorsDoNotUse"
		break
	case StatusTerminateRequested:
		return "ToBeTerminated"
		break
	case StatusDeletingPods:
		return "DeletingPods"
		break
	case StatusDeletingNodes:
		return "DeletingNodes"
		break
	case StatusTerminatingInstances:
		return "TerminatingInstances"
		break
	case StatusTerminated:
		return "Terminated"
		break
	case StatusRejected:
		return "Rejected"
		break
	}
	panic("GetStatusLabel should receive a valid status")
}

// Equal compares if the given status is the same as the current status.
func (s Status) Equal(status int) bool {
	return int(s) == status
}


// ToString returns a string of this status value
func (s Status) ToString() string {
	return string(s)
}

// ErrorStatus are possible status values of SimulationDeployment ErrorStatus field.
type ErrorStatus string

const (
	ErrWhenInitializing ErrorStatus = "InitializationFailed"
	ErrWhenTerminating  ErrorStatus = "TerminationFailed"
	// Set when there was a second error during error handling. Marking for human review
	ErrAdminReview ErrorStatus = "AdminReview"
	// Set when the simulation did not start due to a rejection by the SimService
	ErrRejected ErrorStatus = "Rejected"
	// ErrServerRestart is set by the server initialization process when it finds
	// Simulation Deployments left with intermediate statuses (either starting or terminating).
	// Having this error means that the server suffered a shutdown in the middle of a start
	// or terminate operation.
	ErrServerRestart ErrorStatus = "ServerRestart"
	// Set when there was an error during log upload. Marking for human review
	ErrFailedToUploadLogs ErrorStatus = "FailedToUploadLogs"
)

// ToString returns a string of this status value
func (s ErrorStatus) ToString() string {
	return string(s)
}

// ToStringPtr returns a pointer to string of this status value
func (s ErrorStatus) ToStringPtr() *string {
	str := string(s)
	return &str
}

func (s ErrorStatus) weight() int {
	switch s {
	case ErrWhenInitializing, ErrWhenTerminating:
		return 0
	case ErrRejected:
		return 1
	case ErrServerRestart:
		return 2
	case ErrAdminReview, ErrFailedToUploadLogs:
		return 5
	}
	panic("Invalid value")
}

// isMoreSevere checks if the given error is more severe than the current status.
// Returns true if the given error is more severe than the current status.
func (s ErrorStatus) isMoreSevere(err ErrorStatus) bool {
	return err.weight() > s.weight()
}