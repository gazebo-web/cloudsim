package rules

type Service interface {
	GetByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error)
	GetRemainingSubmissions(owner, circuit string) (*int, error)
}

type service struct {
	repository Repository
}

func (s *service) GetByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error) {
	return s.repository.GetByCircuitAndOwner(ruleType, circuit, owner)
}

func (s *service) GetRemainingSubmissions(owner, circuit string) (*int, error) {
	return s.repository.GetRemainingSubmissions(owner, circuit)
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}