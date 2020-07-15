package tracks

// Service groups a set of methods that have the business logic to perform CRUD operations for a Track.
type Service interface {
	serviceCreate
	serviceRead
	serviceUpdate
	serviceDelete
}

// serviceCreate has the business logic for creating a Track.
type serviceCreate interface {
	Create(track CreateTrackInput) (*CreateTrackOutput, error)
}

// serviceRead has the business logic for reading one or multiple Tracks.
type serviceRead interface {
	Get(name string) (*Track, error)
	GetAll() ([]Track, error)
}

// serviceUpdate has the business logic for updating a Track.
type serviceUpdate interface {
	Update(track UpdateTrackInput) (*UpdateTrackOutput, error)
}

// serviceDelete has the business logic for deleting a Track.
type serviceDelete interface {
	Delete(name string) (*Track, error)
}