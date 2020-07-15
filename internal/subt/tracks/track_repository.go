package tracks

// Repository groups a set of methods to perform CRUD operations in the database for a certain Track.
type Repository interface {
	repositoryCreate
	repositoryRead
	repositoryUpdate
	repositoryDelete
}

// repositoryCreate has a method to Create a track in the database.
type repositoryCreate interface {
	Create(track Track) (*Track, error)
}

// repositoryRead has a method to get one or multiple tracks from the database.
type repositoryRead interface {
	Get(name string)
	GetAll() ([]Track, error)
}

// repositoryUpdate has a method to update a track in the database.
type repositoryUpdate interface {
	Update(name string, track Track) (*Track, error)
}

// repositoryDelete has a method to delete a track from the database.
type repositoryDelete interface {
	Delete(name string) (*Track, error)
}
