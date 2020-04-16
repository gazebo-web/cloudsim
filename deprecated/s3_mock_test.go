package deprecated

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

const (
	S3OpPutObject OpType = "PutObject"
)

// S3Mock is a mock for S3 service.
type S3Mock struct {
	s3iface.S3API
	Mock
	PutObjectFunc func(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

// NewS3Mock creates a new S3Mock.
func NewS3Mock() *S3Mock {
	m := &S3Mock{}
	return m
}

// PutObject API operation for Amazon Simple Storage Service.
// Adds an object to a bucket.
func (m *S3Mock) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	if m.PutObjectFunc != nil {
		return m.PutObjectFunc(input)
	}

	m.Tracker.TrackCall(S3OpPutObject)
	defer m.InvokeCallback(S3OpPutObject, input)

	result := m.GetMockResult(S3OpPutObject)
	// PassThrough is a special value that indicates the non-mocked version of
	// this function should be called
	if result == PassThrough {
		return m.S3API.PutObject(input)
	}
	// If the mock result is an error, return that error
	if err, ok := result.(error); ok {
		return nil, err
	}

	if r, ok := result.(*s3.PutObjectOutput); ok {
		return r, nil
	}
	return nil, nil
}
