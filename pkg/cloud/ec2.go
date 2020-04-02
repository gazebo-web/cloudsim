package cloud

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"reflect"
	"time"
)

// AmazonEC2 wraps the AWS EC2 API.
type AmazonEC2 struct {
	API ec2iface.EC2API
	Retries int
	Delay int
}

// NewAmazonEC2 returns a new AmazonEC2 instance by the given AWS session and configuration.
func NewAmazonEC2(p client.ConfigProvider, cfgs ...*aws.Config) AmazonEC2 {
	var instance AmazonEC2
	if !reflect.ValueOf(p).IsNil() {
		instance.API = ec2.New(p, cfgs...)
	}
	return instance
}

func (ec *AmazonEC2) Launch(ctx context.Context, input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	input.SetDryRun(true)
	for try := 1; try <= ec.Retries; try++ {
		_, err := ec.API.RunInstances(input)
		awsErr, ok := err.(awserr.Error)
		if !ok {
			return nil, err
		}
		if ec.isErrorRetryable(awsErr) {
			logger.Logger(ctx).Info(fmt.Sprintf("[EC2|LAUNCH] Error [%s] while launching nodes on dry mode.\nError: %s\n", awsErr.Code(), awsErr.Message()))
		}
		if ec.isDryRunOperation(awsErr) {
			break
		}
		if try != ec.Retries {
			tools.Sleep(time.Second * time.Duration(ec.Delay))
		}
	}
	input.SetDryRun(false)
	result, err := ec.API.RunInstances(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ec AmazonEC2) isErrorRetryable(err awserr.Error) bool {
	return ec.getRetryableError(err.Code())
}

func (ec AmazonEC2) isDryRunOperation(err awserr.Error) bool {
	return err.Code() == "DryRunOperation"
}

func (ec AmazonEC2) getRetryableError(code string) bool {
	switch code {
	case "RequestLimitExceeded":
	case "InsufficientInstanceCapacity":
		return true
	}
	return false
}