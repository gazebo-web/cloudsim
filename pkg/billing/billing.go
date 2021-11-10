package billing

import (
	"context"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
)

type Config struct {
	CreditsURL  string
	PaymentsURL string
}

type Service interface {
	GetBalance(ctx context.Context, user *users.User) (interface{}, error)
}

type service struct {
}

func (s *service) GetBalance(ctx context.Context, user *users.User) (interface{}, error) {
	panic("implement me")
}

func NewService(cfg Config) Service {
	return &service{}
}
