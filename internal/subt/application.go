package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

type SubT struct {

}

func Register() *application.IApplication {
	var subt application.IApplication
	subt = SubT{}
	return &subt
}
