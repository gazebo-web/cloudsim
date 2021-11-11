package billing

import (
	"context"
	apiCredits "gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	credits "gitlab.com/ignitionrobotics/billing/credits/pkg/client"
	apiPayments "gitlab.com/ignitionrobotics/billing/payments/pkg/api"
	payments "gitlab.com/ignitionrobotics/billing/payments/pkg/client"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"net/url"
	"time"
)

// Service holds methods to interact with different billing services.
type Service interface {
	GetBalance(ctx context.Context, user *users.User) (interface{}, error)
	CreateSession(ctx context.Context, user *users.User, in CreateSessionRequest) (interface{}, error)
}

// service is a Service implementation using the payments and credits V1 API.
type service struct {
	// payments holds a PaymentsV1 client implementation.
	payments payments.Client

	// credits holds a CreditsV1 client implementation.
	credits credits.Client

	// applicationName is the name of the current application that consumes billing services.
	// This value keeps all application transactions in the same context.
	applicationName string
}

// CreateSessionRequest is used to create a new payment session.
type CreateSessionRequest struct {
	// SuccessURL is the url where to redirect users after a payment succeeds.
	SuccessURL string `json:"success_url"`
	// CancelURL is the url where to redirect users after a payment fails.
	CancelURL string `json:"cancel_url"`
}

// CreateSession creates a new payment session.
func (s *service) CreateSession(ctx context.Context, user *users.User, in CreateSessionRequest) (interface{}, error) {
	return s.payments.CreateSession(ctx, apiPayments.CreateSessionRequest{
		Service:     "stripe",
		SuccessURL:  in.SuccessURL,
		CancelURL:   in.CancelURL,
		Handle:      *user.Username,
		Application: s.applicationName,
	})
}

// GetBalance returns the credits balance of the given user.
func (s *service) GetBalance(ctx context.Context, user *users.User) (interface{}, error) {
	return s.credits.GetBalance(ctx, apiCredits.GetBalanceRequest{
		Handle:      *user.Username,
		Application: s.applicationName,
	})
}

// Config is used to configure new service implementations
type Config struct {
	CreditsURL      string
	PaymentsURL     string
	ApplicationName string
	Timeout         time.Duration
}

// NewService initializes a new Service implementation using the given config.
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
