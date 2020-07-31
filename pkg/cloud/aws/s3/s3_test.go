package s3

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestUpload_OK(t *testing.T) {
	output := s3.PutObjectOutput{}
	err := runUploadTest(&output, nil)
	assert.NoError(t, err)
}

func TestUpload_Failed(t *testing.T) {
	err := runUploadTest(nil, errors.New("s3-test"))
	assert.Error(t, err)
}

func TestGetURL_OK(t *testing.T) {
	bucket := "test"
	filepath := "tmp/test/log.txt"
	u := &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.s3.amazonws.com", bucket),
		Path:   filepath,
	}
	result, err := runGetURLTest(bucket, filepath, &request.Request{
		Operation: &request.Operation{
			BeforePresignFn: nil,
		},

		HTTPRequest: &http.Request{
			Method: "GET",
			URL:    u,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, u.String(), result)
}

func TestGetURL_Failed(t *testing.T) {
	bucket := "test"
	filepath := "/tmp/test/log.txt"
	u := &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.s3.amazonws.com", bucket),
		Path:   filepath,
	}
	result, err := runGetURLTest(bucket, filepath, &request.Request{
		Operation: &request.Operation{
			BeforePresignFn: nil,
		},
		Error: errors.New(request.ErrCodeInvalidPresignExpire),
		HTTPRequest: &http.Request{
			Method: "GET",
			URL:    u,
		},
	})
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

func runUploadTest(desiredOutput *s3.PutObjectOutput, err error) error {
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

func runGetURLTest(bucket, filepath string, req *request.Request) (string, error) {
	api := &s3api{
		Mock: new(mock.Mock),
	}

	expiresIn := time.Second

	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &filepath,
	}

	api.On("GetObjectRequest", input).Return(req)

	s := NewStorage(api)

	return s.GetURL(bucket, filepath, expiresIn)
}

type s3api struct {
	*mock.Mock
	s3iface.S3API
}

func (s *s3api) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := s.Called(input)
	output := args.Get(0).(*s3.PutObjectOutput)
	err := args.Error(1)
	return output, err
}

func (s *s3api) GetObjectRequest(input *s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	args := s.Called(input)
	req := args.Get(0).(*request.Request)
	return req, nil
}
