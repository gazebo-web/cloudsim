package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

func RegisterApplications(p *platform.Platform) {
	p.Applications = append(p.Applications, subt.Register())
}