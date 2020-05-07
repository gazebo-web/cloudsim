package cloud

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"
)

type AmazonS3Mock struct {
	*mock.Mock
}

func NewAmazonS3Mock() *AmazonS3Mock {
	var ec2 *AmazonS3Mock
	ec2 = &AmazonS3Mock{
		Mock: new(mock.Mock),
	}
	return ec2
}

func (s *AmazonS3Mock) GetAddress(bucket string, key string) string {
	args := s.Called(bucket, key)
	return args.String(0)
}

func (s *AmazonS3Mock) Upload(bucket string, key string, file []byte) (*s3.PutObjectOutput, error) {
	args := s.Called(bucket, key, file)
	output := args.Get(0).(*s3.PutObjectOutput)
	return output, args.Error(1)
}

func (s *AmazonS3Mock) GetLogKey(groupID string, owner string) string {
	args := s.Called(groupID, owner)
	return args.String(0)
}