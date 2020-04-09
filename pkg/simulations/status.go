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