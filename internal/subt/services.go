package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

// Services extends a generic application services interface to add SubT services.
type Services interface {
	application.Services
	Tracks() tracks.Service
}
