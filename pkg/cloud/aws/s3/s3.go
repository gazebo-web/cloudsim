package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"time"
)

// storage is a cloud.Storage implementation.
type storage struct {
	API s3iface.S3API
}

// Upload uploads a file to the cloud storage.
func (s storage) Upload(input cloud.UploadInput) error {
	panic("implement me")
}

// GetURL returns an URL to access the given bucket with the given key.
func (s storage) GetURL(bucket string, key string, expiresIn time.Duration) string {
	panic("implement me")
}

// NewStorage initializes a new cloud.Storage implementation using s3.
func NewStorage(api s3iface.S3API) cloud.Storage {
	return &storage{
		API: api,
	}
}
