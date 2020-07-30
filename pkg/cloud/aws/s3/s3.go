package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"time"
)

type storage struct {
	API s3iface.S3API
}

func (s storage) Upload(input cloud.UploadInput) error {
	panic("implement me")
}

func (s storage) GetURL(bucket string, key string, expireIn time.Duration) string {
	panic("implement me")
}

func NewStorage(api s3iface.S3API) cloud.Storage {
	return &storage{
		API: api,
	}
}
