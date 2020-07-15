package tracks

// Service groups a set of methods that have the business logic to perform CRUD operations for a Circuit.
type Service interface {
	serviceCreate
}

// serviceCreate has the business logic for creating a Track.
type serviceCreate interface {
	Create(track CreateTrackInput) (*CreateTrackOutput, error)
}
