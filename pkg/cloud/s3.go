package cloud

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"net/http"
	"path/filepath"
	"reflect"
)

// AmazonS3 wraps the AWS S3 API.
type AmazonS3 struct {
	API s3iface.S3API
}

// NewAmazonS3 returns a new AmazonS3 instance by the given AWS session and configuration.
func NewAmazonS3(p client.ConfigProvider, cfgs ...*aws.Config) *AmazonS3 {
	var instance AmazonS3
	if !reflect.ValueOf(p).IsNil() {
		instance.API = s3.New(p, cfgs...)
	}
	return &instance
}

func (s AmazonS3) Address(bucket string, key string) string {
	return fmt.Sprintf("s3://%s", filepath.Join(bucket, key))
}

func (s AmazonS3) Upload(bucket string, key string, file []byte) (*s3.PutObjectOutput, error) {
	return s.API.PutObject(&s3.PutObjectInput{
		Bucket:               &bucket,
		Key:                  &key,
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(file),
		ContentLength:        aws.Int64(int64(len(file))),
		ContentType:          aws.String(http.DetectContentType(file)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
}