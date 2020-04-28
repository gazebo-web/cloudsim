package circuits

type IService interface {
	GetByName(name string) (*Circuit, error)
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	var s IService
	s = &Service{
		repository: repository,
	}
	return s
}

func (s *Service) GetByName(name string) (*Circuit, error) {
	panic("implement me")
}
