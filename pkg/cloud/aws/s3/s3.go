package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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
	_, err := s.API.PutObject(&s3.PutObjectInput{
		Bucket:               &input.Bucket,
		Key:                  &input.Key,
		ACL:                  aws.String("private"),
		Body:                 input.File,
		ContentLength:        aws.Int64(input.ContentLength),
		ContentType:          aws.String(input.ContentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

// GetURL returns an URL to access the given bucket with the given key.
func (s storage) GetURL(bucket string, key string, expiresIn time.Duration) (string, error) {
	req, _ := s.API.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	u, err := req.Presign(expiresIn)
	if err != nil {
		return "", err
	}
	return u, nil
}

// NewStorage initializes a new cloud.Storage implementation using s3.
func NewStorage(api s3iface.S3API) cloud.Storage {
	return &storage{
		API: api,
	}
}
