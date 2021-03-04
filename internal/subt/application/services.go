package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/summaries"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

// Services extends a generic application services interface to add SubT services.
type Services interface {
	application.Services
	Tracks() tracks.Service
	Summaries() summaries.Service
}

// services is a Services implementation.
type services struct {
	application.Services
	tracks    tracks.Service
	summaries summaries.Service
}

func (s *services) Summaries() summaries.Service {
	return s.summaries
}

// Tracks returns a Track service.
func (s *services) Tracks() tracks.Service {
	return s.tracks
}

// NewServices initializes a new Services implementation using a base generic service.
func NewServices(base application.Services, tracks tracks.Service, summaries summaries.Service) Services {
	return &services{
		Services:  base,
		tracks:    tracks,
		summaries: summaries,
	}
}
