package state

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"

// Platform exposes a method to access the platform.
type Platform interface {
	Platform() platform.Platform
}
