package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/statistics"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

// Services extends a generic application services interface to add SubT services.
type Services interface {
	application.Services
	Tracks() tracks.Service
	Statistics() statistics.Service
}

// services is a Services implementation.
type services struct {
	application.Services
	tracks     tracks.Service
	statistics statistics.Service
}

func (s *services) Statistics() statistics.Service {
	return s.statistics
}

// Tracks returns a Track service.
func (s *services) Tracks() tracks.Service {
	return s.tracks
}

// NewServices initializes a new Services implementation using a base generic service.
func NewServices(base application.Services, tracks tracks.Service, statistics statistics.Service) Services {
	return &services{
		Services:   base,
		tracks:     tracks,
		statistics: statistics,
	}
}
