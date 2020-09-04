package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

type Services interface {
	application.Services
	Tracks() tracks.Service
}
