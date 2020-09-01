package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
	"time"
)

// NewAPI returns an S3 client from the given config provider.
func NewAPI(config client.ConfigProvider) s3iface.S3API {
	return s3.New(config)
}

// storage is a cloud.Storage implementation.
type storage struct {
	API    s3iface.S3API
	Logger ign.Logger
}

// Upload uploads a file to the cloud storage.
func (s storage) Upload(input cloud.UploadInput) error {
	s.Logger.Debug(fmt.Sprintf("Upload input: %+v", input))
	if input.File == nil {
		return cloud.ErrMissingFile
	}

	var bslice []byte
	_, err := input.File.Read(bslice)
	if err != nil {
		s.Logger.Debug(fmt.Sprintf("Reading file failed. Error: %s", err))
		return err
	}

	if http.DetectContentType(bslice) != input.ContentType {
		s.Logger.Debug(fmt.Sprintf("Invalid content type. Actual: %s. Expected: %s.", http.DetectContentType(bslice), input.ContentType))
		return cloud.ErrBadContentType
	}

	_, err = s.API.PutObject(&s3.PutObjectInput{
		Bucket:               &input.Bucket,
		Key:                  &input.Key,
		ACL:                  aws.String("private"),
		Body:                 input.File,
		ContentLength:        aws.Int64(input.ContentLength),
		ContentType:          aws.String(input.ContentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		s.Logger.Debug(fmt.Sprintf("Uploading object failed. Error: %s", err))
		return err
	}
	s.Logger.Debug(fmt.Sprintf("Uploading file to bucket [%s] succeeded.", input.Bucket))
	return nil
}

// GetURL returns an URL to access the given bucket with the given key.
func (s storage) GetURL(bucket string, key string, expiresIn time.Duration) (string, error) {
	s.Logger.Debug(fmt.Sprintf("Getting URL for bucket [%s] and key [%s]. Expiration in: %s", bucket, key, expiresIn.String()))
	req, _ := s.API.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	u, err := req.Presign(expiresIn)
	if err != nil {
		s.Logger.Debug(fmt.Sprintf("Presigning URL failed. Error: %s", err))
		return "", err
	}
	s.Logger.Debug(fmt.Sprintf("Getting URL succeeded. URL: %s", u))
	return u, nil
}

// NewStorage initializes a new cloud.Storage implementation using s3.
func NewStorage(api s3iface.S3API, logger ign.Logger) cloud.Storage {
	return &storage{
		API:    api,
		Logger: logger,
	}
}
