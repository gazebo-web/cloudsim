package workers

type LaunchDTO struct {
	GroupID string
	Action interface{}
}

type TerminateDTO struct {
	GroupID string
	Action interface{}
}
