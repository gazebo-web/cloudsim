package billing

import (
	"context"
	apiCredits "gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	credits "gitlab.com/ignitionrobotics/billing/credits/pkg/client"
	apiPayments "gitlab.com/ignitionrobotics/billing/payments/pkg/api"
	payments "gitlab.com/ignitionrobotics/billing/payments/pkg/client"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/url"
	"time"
)

// CreateSessionRequest is used to create a new payment session.
type CreateSessionRequest struct {
	// SuccessURL is the url where to redirect users after a payment succeeds.
	SuccessURL string `json:"success_url"`

	// CancelURL is the url where to redirect users after a payment fails.
	CancelURL string `json:"cancel_url"`

	// Handle is the username of the user that is starting the payment session.
	Handle string `json:"-"`
}

// CreateSessionResponse is the response from calling Service.CreateSession.
type CreateSessionResponse apiPayments.CreateSessionResponse

// GetBalanceResponse is the response from calling Service.GetBalance.
type GetBalanceResponse apiCredits.GetBalanceResponse

// Service holds methods to interact with different billing services.
type Service interface {
	// GetBalance returns the credits balance of the given user.
	GetBalance(ctx context.Context, user *users.User) (GetBalanceResponse, error)
	// CreateSession creates a new payment session.
	CreateSession(ctx context.Context, in CreateSessionRequest) (CreateSessionResponse, error)
	// IsEnabled returns true when this service is enabled.
	IsEnabled() bool
	// SubtractCredits subtracts the credits from the given user for the amount of time the given simulation has been running.
	SubtractCredits(ctx context.Context, user *users.User, sim simulations.Simulation) error
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

	// logger is used to log relevant information in different methods.
	logger ign.Logger

	// enabled is set to true when this service is enabled.
	enabled bool
}

// SubtractCredits subtracts the credits from the given user for the amount of time the given simulation has been running.
func (s *service) SubtractCredits(ctx context.Context, user *users.User, sim simulations.Simulation) error {
	rate := sim.GetRate()

	price, err := sim.ApplyRate()
	if err != nil {
		return err
	}

	_, err = s.credits.DecreaseCredits(ctx, apiCredits.DecreaseCreditsRequest{
		Transaction: apiCredits.Transaction{
			Handle:      *user.Username,
			Amount:      price,
			Currency:    rate.Currency,
			Application: s.applicationName,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// IsEnabled returns true when this service is enabled.
func (s *service) IsEnabled() bool {
	return s.enabled
}

// CreateSession creates a new payment session.
func (s *service) CreateSession(ctx context.Context, req CreateSessionRequest) (CreateSessionResponse, error) {
	s.logger.Debug("Creating session:", req)
	res, err := s.payments.CreateSession(ctx, apiPayments.CreateSessionRequest{
		Service:     "stripe",
		SuccessURL:  req.SuccessURL,
		CancelURL:   req.CancelURL,
		Handle:      req.Handle,
		Application: s.applicationName,
	})
	if err != nil {
		return CreateSessionResponse{}, err
	}
	return CreateSessionResponse(res), nil
}

// GetBalance returns the credits balance of the given user.
func (s *service) GetBalance(ctx context.Context, user *users.User) (GetBalanceResponse, error) {
	s.logger.Debug("Getting balance for user:", *user.Username)
	res, err := s.credits.GetBalance(ctx, apiCredits.GetBalanceRequest{
		Handle:      *user.Username,
		Application: s.applicationName,
	})
	if err != nil {
		s.logger.Debug("Failed to get balance:", err)
		return GetBalanceResponse{}, err
	}
	return GetBalanceResponse(res), nil
}

// Config is used to configure new service implementations
type Config struct {
	// CreditsURL contains the URL of the Credits API.
	CreditsURL string
	// PaymentsURL contains the URL of the Payments API.
	PaymentsURL string
	// ApplicationName contains the unique name for this application. Used to track billing operations and keep them
	// in the same application context.
	ApplicationName string
	// Timeout is the amount of time a client can wait until a timeout occurs.
	Timeout time.Duration
	// Enabled is set to true if the service should be enabled.
	Enabled bool
}

// NewService initializes a new Service implementation using the given config.
func NewService(cfg Config, logger ign.Logger) (Service, error) {
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
		logger:          logger,
		enabled:         cfg.Enabled,
	}, nil
}
