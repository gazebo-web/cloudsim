package state

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/application"
)

// ServicesGetter exposes a method to access the application services.
type ServicesGetter interface {
	Services() application.Services
}
