package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

// Services extends a generic application services interface to add SubT services.
type Services interface {
	application.Services
	Tracks() tracks.Service
}

// services is a Services implementation.
type services struct {
	application.Services
	tracks tracks.Service
}

// Tracks returns a Track service.
func (s *services) Tracks() tracks.Service {
	return s.tracks
}

// NewServices initializes a new Services implementation using a base generic service.
func NewServices(base application.Services, tracks tracks.Service) Services {
	return &services{
		Services: base,
		tracks:   tracks,
	}
}
