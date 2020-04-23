package workers

// LaunchInput
type LaunchInput struct {
	GroupID string
	Action interface{}
}

// TerminateInput
type TerminateInput struct {
	GroupID string
	Action interface{}
}
