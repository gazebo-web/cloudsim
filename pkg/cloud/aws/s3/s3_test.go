package s3

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"testing"
)

type s3api struct {
	*mock.Mock
	s3iface.S3API
}

func runStorageTest(desiredOutput *s3.PutObjectOutput, err error) error {
	api := &s3api{
		Mock: new(mock.Mock),
	}

	bucket := "test"
	key := "1234"
	body := []byte("test")
	file := bytes.NewReader(body)
	fileSize := int64(len(body))
	contentType := "type-test"

	input := &s3.PutObjectInput{
		Bucket:               &bucket,
		Key:                  &key,
		ACL:                  aws.String("private"),
		Body:                 file,
		ContentLength:        &fileSize,
		ContentType:          &contentType,
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	}

	api.On("PutObject", input).Return(desiredOutput, err)

	s := NewStorage(api)

	return s.Upload(cloud.UploadInput{
		Bucket:        bucket,
		Key:           key,
		File:          file,
		ContentLength: fileSize,
		ContentType:   contentType,
	})
}

func (s *s3api) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := s.Called(input)
	output := args.Get(0).(*s3.PutObjectOutput)
	err := args.Error(1)
	return output, err
}

func TestUpload_OK(t *testing.T) {
	output := s3.PutObjectOutput{}
	err := runStorageTest(&output, nil)
	assert.NoError(t, err)
}

func TestUpload_Failed(t *testing.T) {
	err := runStorageTest(nil, errors.New("s3-test"))
	assert.Error(t, err)
}
