package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
	"net/http"
	"path/filepath"
	"time"
)

// NewAPI returns an S3 client from the given config provider.
func NewAPI(config client.ConfigProvider) s3iface.S3API {
	return s3.New(config)
}

// s3Storage is a storage.Storage implementation.
type s3Storage struct {
	API    s3iface.S3API
	Logger ign.Logger
}

func (s s3Storage) PrepareAddress(bucket, key string) string {
	return fmt.Sprintf("s3://%s", filepath.Join(bucket, key))
}

// Upload uploads a file to the cloud storage.
func (s s3Storage) Upload(input storage.UploadInput) error {
	s.Logger.Debug(fmt.Sprintf("Upload input: %+v", input))
	if input.File == nil {
		return storage.ErrMissingFile
	}

	var bslice []byte
	_, err := input.File.Read(bslice)
	if err != nil {
		s.Logger.Debug(fmt.Sprintf("Reading file failed. Error: %s", err))
		return err
	}

	if http.DetectContentType(bslice) != input.ContentType {
		s.Logger.Debug(fmt.Sprintf("Invalid content type. Actual: %s. Expected: %s.", http.DetectContentType(bslice), input.ContentType))
		return storage.ErrBadContentType
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
func (s s3Storage) GetURL(bucket string, key string, expiresIn time.Duration) (string, error) {
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

// NewStorage initializes a new storage.Storage implementation using s3.
func NewStorage(api s3iface.S3API, logger ign.Logger) storage.Storage {
	return &s3Storage{
		API:    api,
		Logger: logger,
	}
}
