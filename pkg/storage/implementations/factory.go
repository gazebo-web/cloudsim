package implementations

import (
	factorymap "github.com/gazebo-web/cloudsim/v4/pkg/factory/map"
	s3factory "github.com/gazebo-web/cloudsim/v4/pkg/storage/implementations/s3/factory"
)

const (
	// S3 is the S3 implementation factory identifier.
	S3 = "s3"
)

// Factory provides a factory to create Storage implementations.
var Factory = factorymap.Map{
	S3: s3factory.NewFunc,
}
