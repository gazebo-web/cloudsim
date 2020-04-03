package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"reflect"
	"sync"
)

// AmazonEC2 wraps the AWS EC2 API.
type AmazonEC2 struct {
	API ec2iface.EC2API
	Retries int
	lockLaunch sync.Mutex
}

// NewAmazonEC2 returns a new AmazonEC2 instance by the given AWS session and configuration.
func NewAmazonEC2(p client.ConfigProvider, cfgs ...*aws.Config) *AmazonEC2 {
	var instance AmazonEC2
	if !reflect.ValueOf(p).IsNil() {
		instance.API = ec2.New(p, cfgs...)
	}
	return &instance
}