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
	Create(track CreateTrackInput) (*Track, error)
}

// serviceRead has the business logic for reading one or multiple Tracks.
type serviceRead interface {
	Get(name string) (*Track, error)
	GetAll() ([]Track, error)
}

// serviceUpdate has the business logic for updating a Track.
type serviceUpdate interface {
	Update(track UpdateTrackInput) (*Track, error)
}

// serviceDelete has the business logic for deleting a Track.
type serviceDelete interface {
	Delete(name string) (*Track, error)
}

type service struct {
	repository Repository
}

func (s service) Create(track CreateTrackInput) (*Track, error) {
	panic("implement me")
}

func (s service) Get(name string) (*Track, error) {
	panic("implement me")
}

func (s service) GetAll() ([]Track, error) {
	panic("implement me")
}

func (s service) Update(track UpdateTrackInput) (*Track, error) {
	panic("implement me")
}

func (s service) Delete(name string) (*Track, error) {
	panic("implement me")
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
