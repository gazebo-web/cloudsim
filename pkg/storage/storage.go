package storage

import (
	"bytes"
	"errors"
	"time"
)

var (
	// ErrMissingFile is returned when Upload gets called and no file is provided.
	ErrMissingFile = errors.New("missing file")
	// ErrBadContentType is returned when Upload gets called with an ivnalid content type.
	ErrBadContentType = errors.New("bad content type")
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
	// Upload uploads a file to a cloud storage.
	Upload(input UploadInput) error
	// GetURL returns the URL of the given bucket and key from a cloud storage.
	GetURL(bucket string, key string, expiresIn time.Duration) (string, error)
	// PrepareAddress returns the address for the given bucket with the given key.
	PrepareAddress(bucket, key string) string
}
