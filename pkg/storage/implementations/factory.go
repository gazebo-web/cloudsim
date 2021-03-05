package implementations

import (
	factorymap "gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory/map"
	s3factory "gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage/implementations/s3/factory"
)

const (
	// S3 is the S3 implementation factory identifier.
	S3 = "s3"
)

// Factory provides a factory to create Storage implementations.
var Factory = factorymap.Map{
	S3: s3factory.NewFunc,
}
