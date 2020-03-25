package application

import "gitlab.com/ignitionrobotics/web/ign-go"

type IApplication interface {
	RegisterRoutes() ign.Routes
}
