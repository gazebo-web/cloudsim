package workers

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"

// LaunchInput is the input needed to launch a simulation.
type LaunchInput struct {
	GroupID string
	Action  *actions.Action
}

// TerminateInput is the input needed to terminate a simulation.
type TerminateInput struct {
	GroupID string
	Action  *actions.Action
}
