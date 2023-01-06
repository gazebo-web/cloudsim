package state

import "github.com/gazebo-web/cloudsim/v4/pkg/platform"

// PlatformGetter exposes a method to access the platform.
type PlatformGetter interface {
	Platform() platform.Platform
}
