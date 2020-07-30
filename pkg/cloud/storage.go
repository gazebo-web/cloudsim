package cloud

import (
	"bytes"
	"time"
)

// UploadInput is the input for the Storage.Upload operation.
// It will be used to upload a file to a certain bucket.
type UploadInput struct {
	Bucket        string
	Key           string
	File          *bytes.Reader
	ContentLength int64
	ContentType   string
}

// Storage groups a set of methods to interact with a Cloud Storage.
type Storage interface {
	Upload(input UploadInput) error
	GetURL(bucket string, key string, expiresIn time.Duration) (string, error)
}
