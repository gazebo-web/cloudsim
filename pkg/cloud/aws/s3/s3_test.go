package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(s3StorageTestSuite))
}

type s3test struct {
	Backend *s3mem.Backend
	Faker   *gofakes3.GoFakeS3
	Server  *httptest.Server
	Config  *aws.Config
	Session *session.Session
}

type s3StorageTestSuite struct {
	suite.Suite
	s3      *s3test
	storage cloud.Storage
}

func (s *s3StorageTestSuite) SetupTest() {
	backend := s3mem.New()
	faker := gofakes3.New(backend)
	server := httptest.NewServer(faker.Server())
	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
		Endpoint:         aws.String(server.URL),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	s.s3 = &s3test{
		Backend: backend,
		Faker:   faker,
		Server:  server,
		Config:  config,
		Session: session.Must(session.NewSession(config)),
	}
}

func (s *s3StorageTestSuite) AfterTest() {
	s.s3.Server.Close()
}

func (s *s3StorageTestSuite) TestNewStorage() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)
	obj, ok := st.(*storage)
	s.True(ok)
	s.NotNil(obj.API)
}

func (s *s3StorageTestSuite) TestUpload_OK() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)
	_, err := api.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("bucket")})
	s.NoError(err)

	bslice := []byte("test")
	file := bytes.NewReader(bslice)
	err = st.Upload(cloud.UploadInput{
		Bucket:        "bucket",
		Key:           "key",
		File:          file,
		ContentLength: int64(len(bslice)),
		ContentType:   http.DetectContentType(bslice),
	})
	s.NoError(err)
}

func (s *s3StorageTestSuite) TestUpload_BucketDoesntExist() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)

	bslice := []byte("test")
	file := bytes.NewReader(bslice)
	err := st.Upload(cloud.UploadInput{
		Bucket:        "bucket",
		Key:           "key",
		File:          file,
		ContentLength: int64(len(bslice)),
		ContentType:   http.DetectContentType(bslice),
	})
	s.Error(err)
}

func (s *s3StorageTestSuite) TestUpload_MissingBucketName() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)

	_, err := api.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("bucket")})
	s.NoError(err)

	bslice := []byte("test")
	file := bytes.NewReader(bslice)
	err = st.Upload(cloud.UploadInput{
		Bucket:        "",
		Key:           "key",
		File:          file,
		ContentLength: int64(len(bslice)),
		ContentType:   http.DetectContentType(bslice),
	})
	s.Error(err)
}

func (s *s3StorageTestSuite) TestUpload_MissingKey() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)

	_, err := api.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("bucket")})
	s.NoError(err)

	bslice := []byte("test")
	file := bytes.NewReader(bslice)
	err = st.Upload(cloud.UploadInput{
		Bucket:        "bucket",
		Key:           "",
		File:          file,
		ContentLength: int64(len(bslice)),
		ContentType:   http.DetectContentType(bslice),
	})
	s.Error(err)
}

func (s *s3StorageTestSuite) TestUpload_MissingFile() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)

	_, err := api.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("bucket")})
	s.NoError(err)

	bslice := []byte("test")
	err = st.Upload(cloud.UploadInput{
		Bucket:        "bucket",
		Key:           "key",
		File:          nil,
		ContentLength: int64(len(bslice)),
		ContentType:   http.DetectContentType(bslice),
	})
	s.Error(err)
	s.Equal(cloud.ErrMissingFile, err)
}

func (s *s3StorageTestSuite) TestUpload_FileLengthMismatch() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)

	_, err := api.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("bucket")})
	s.NoError(err)

	bslice := []byte("test")
	file := bytes.NewReader(bslice)
	err = st.Upload(cloud.UploadInput{
		Bucket:        "bucket",
		Key:           "key",
		File:          file,
		ContentLength: int64(len(bslice) + 123),
		ContentType:   http.DetectContentType(bslice),
	})
	s.Error(err)
}

func (s *s3StorageTestSuite) TestUpload_BadContentType() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)

	_, err := api.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("bucket")})
	s.NoError(err)

	bslice := []byte("test")
	file := bytes.NewReader(bslice)
	err = st.Upload(cloud.UploadInput{
		Bucket:        "bucket",
		Key:           "key",
		File:          file,
		ContentLength: int64(len(bslice)),
		ContentType:   "test",
	})
	s.Error(err)
	s.Equal(cloud.ErrBadContentType, err)
}

func (s *s3StorageTestSuite) TestGetURL() {
	api := s3.New(s.s3.Session)
	st := NewStorage(api)

	_, err := api.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("bucket")})
	s.NoError(err)

	bslice := []byte("test")
	file := bytes.NewReader(bslice)
	err = st.Upload(cloud.UploadInput{
		Bucket:        "bucket",
		Key:           "key",
		File:          file,
		ContentLength: int64(len(bslice)),
		ContentType:   http.DetectContentType(bslice),
	})
	s.NoError(err)

	u, err := st.GetURL("bucket", "key", 5*time.Minute)
	var expectedType string
	s.NoError(err)
	s.IsType(expectedType, u)

	_, err = url.Parse(u)
	s.NoError(err)
}
