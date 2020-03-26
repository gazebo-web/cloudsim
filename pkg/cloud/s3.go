package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"reflect"
)

// AmazonS3 wraps the AWS S3 API.
type AmazonS3 struct {
	API s3iface.S3API
}

// NewAmazonS3 returns a new AmazonS3 instance by the given AWS session and configuration.
func NewAmazonS3(p client.ConfigProvider, cfgs ...*aws.Config) AmazonS3 {
	var instance AmazonS3
	if !reflect.ValueOf(p).IsNil() {
		instance.API = s3.New(p, cfgs...)
	}
	return instance
}