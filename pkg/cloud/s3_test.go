package cloud

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAmazonS3(t *testing.T) {
	s := session.Must(session.NewSession())
	s3 := NewAmazonS3(s)
	assert.NotNil(t, s3)
}

func TestAmazonS3_GetAddress(t *testing.T) {
	s := session.Must(session.NewSession())
	s3 := NewAmazonS3(s)
	bucket := "bucket_test"
	key := "key_test"
	expected := fmt.Sprintf("s3://%s/%s", bucket, key)
	addr := s3.GetAddress(bucket, key)
	assert.Equal(t, expected, addr)
}

func TestAmazonS3_GetLogKey(t *testing.T) {
	s := session.Must(session.NewSession())
	s3 := NewAmazonS3(s)
	groupID := "test-test-test-test"
	owner := "Open Robotics"
	escaped := "Open%20Robotics"
	expected := fmt.Sprintf("/gz-logs/%s/%s/", escaped, groupID)
	logKey := s3.GetLogKey(groupID, owner)
	assert.Equal(t, expected, logKey)
}
