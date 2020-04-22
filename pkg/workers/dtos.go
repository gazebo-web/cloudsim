package workers

// LaunchDTO
type LaunchDTO struct {
	GroupID string
	Action interface{}
}

// TerminateDTO
type TerminateDTO struct {
	GroupID string
	Action interface{}
}
