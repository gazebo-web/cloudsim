package billing

import (
	"context"
	apiCredits "gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	credits "gitlab.com/ignitionrobotics/billing/credits/pkg/client"
	payments "gitlab.com/ignitionrobotics/billing/payments/pkg/client"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"net/url"
	"time"
)

type Config struct {
	CreditsURL      string
	PaymentsURL     string
	ApplicationName string
	Timeout         time.Duration
}

type Service interface {
	GetBalance(ctx context.Context, user *users.User) (interface{}, error)
}

type service struct {
	payments        payments.Client
	credits         credits.Client
	applicationName string
}

func (s *service) GetBalance(ctx context.Context, user *users.User) (interface{}, error) {
	return s.credits.GetBalance(ctx, apiCredits.GetBalanceRequest{
		Handle:      *user.Username,
		Application: s.applicationName,
	})
}

func NewService(cfg Config) (Service, error) {
	u, err := url.Parse(cfg.PaymentsURL)
	if err != nil {
		return nil, err
	}
	p := payments.NewPaymentsClientV1(u, cfg.Timeout)

	u, err = url.Parse(cfg.CreditsURL)
	if err != nil {
		return nil, err
	}
	c := credits.NewCreditsClientV1(u, cfg.Timeout)
	return &service{
		payments:        p,
		credits:         c,
		applicationName: cfg.ApplicationName,
	}, nil
}
