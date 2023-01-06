package state

import "github.com/gazebo-web/cloudsim/pkg/platform"

// PlatformGetter exposes a method to access the platform.
type PlatformGetter interface {
	Platform() platform.Platform
}
