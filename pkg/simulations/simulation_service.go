package simulations

// IService
type IService interface {
	GetRepository() IRepository
	SetRepository(repository IRepository)
	Get(groupID string) (*Simulation, error)
	GetAll() []Simulation
	GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(statusFrom, statusTo Status) (*Simulations, error)
	GetAllParentsWithErrors(statusFrom, statusTo Status, errors []ErrorStatus) (*Simulations, error)
	GetParent(groupID string) (*Simulation, error)
	Update(groupID string, simulation Simulation) (*Simulation, error)
	UpdateParentFromChildren(parent *Simulation) (*Simulation, error)
}

// Service
type Service struct {
	repository IRepository
}

// NewService
func NewService(repository IRepository) IService {
	var s IService
	s = &Service{repository: repository}
	return s
}

// GetRepository
func (s *Service) GetRepository() IRepository {
	return s.repository
}

// SetRepository
func (s *Service) SetRepository(repository IRepository) {
	s.repository = repository
}

// Get
func (s *Service) Get(groupID string) (*Simulation, error) {
	panic("Not implemented")
}

// GetAll
func (s *Service) GetAll() []Simulation {
	panic("Not implemented")

}

// GetAllByOwner
func (s *Service) GetAllByOwner(owner string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")

}

// GetChildren
func (s *Service) GetChildren(groupID string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

// GetAllParents
func (s *Service) GetAllParents(statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

// GetAllParentsWithErrors
func (s * Service) GetAllParentsWithErrors(statusFrom, statusTo Status, errors []ErrorStatus) (*Simulations, error) {
	panic("Not implemented")
}

// GetParent
func (s *Service) GetParent(groupID string) (*Simulation, error) {
	panic("implement me")
}

// UpdateParentFromChildren
func (s *Service) UpdateParentFromChildren(parent *Simulation) (*Simulation, error) {
	panic("implement me")
}

// Update
func (s *Service) Update(groupID string, simulation Simulation) (*Simulation, error) {
	sim, err := s.repository.Update(groupID, simulation)
	if err != nil {
		return nil, err
	}
	return sim, nil
}