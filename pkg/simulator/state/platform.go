package state

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"

// PlatformGetter exposes a method to access the platform.
type PlatformGetter interface {
	Platform() platform.Platform
}
