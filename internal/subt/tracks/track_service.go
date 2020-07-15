package tracks

type Service interface {
	serviceCreate
}

type serviceCreate interface {
	Create(track CreateTrackInput) (*CreateTrackOutput, error)
}
