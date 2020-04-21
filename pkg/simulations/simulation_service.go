package simulations

type IService interface {
	GetRepository() IRepository
	SetRepository(repository IRepository)
	Get(groupID string) (*Simulation, error)
	GetAll() []Simulation
	GetAllByOwner(owner string, application string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, application string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(application string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParentsWithErrors(application string, statusFrom, statusTo Status, errors []ErrorStatus) (*Simulations, error)
	Update(groupID string, simulation Simulation) (*Simulation, error)
	UpdateParentFromChildren(parent *Simulation) (*Simulation, error)
	GetParent(application string, groupID string) (*Simulation, error)
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	var s IService
	s = &Service{repository: repository}
	return s
}

func (s *Service) GetRepository() IRepository {
	return s.repository
}

func (s *Service) SetRepository(repository IRepository) {
	s.repository = repository
}

func (s *Service) Update(groupID string, simulation Simulation) (*Simulation, error) {
	sim, err := s.repository.Update(groupID, simulation)
	if err != nil {
		return nil, err
	}
	return sim, nil
}

func (s *Service) Get(groupID string) (*Simulation, error) {
	panic("Not implemented")
}

func (s *Service) GetAll() []Simulation {
	panic("Not implemented")

}

func (s *Service) GetAllByOwner(owner string, application string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")

}

func (s *Service) GetChildren(groupID string, application string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

func (s *Service) GetAllParents(application string, statusFrom, statusTo Status) (*Simulations, error) {
	panic("Not implemented")
}

func (s * Service) GetAllParentsWithErrors(application string, statusFrom, statusTo Status, errors []ErrorStatus) (*Simulations, error) {
	panic("Not implemented")
}

func (s *Service) UpdateParentFromChildren(parent *Simulation) (*Simulation, error) {
	panic("implement me")
}

func (s *Service) GetParent(application string, groupID string) (*Simulation, error) {
	panic("implement me")
}