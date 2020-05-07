package workers

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"

// LaunchInput
type LaunchInput struct {
	GroupID string
	Action actions.Action
}

// TerminateInput
type TerminateInput struct {
	GroupID string
	Action actions.Action
}
