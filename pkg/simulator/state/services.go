package state

import (
	"github.com/gazebo-web/cloudsim/pkg/application"
)

// ServicesGetter exposes a method to access the application services.
type ServicesGetter interface {
	Services() application.Services
}
