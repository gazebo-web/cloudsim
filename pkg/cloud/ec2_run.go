package cloud

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"time"
)

type RunInstanceConfig struct {
	DryRun           bool
	KeyName          string
	MinCount         int64
	MaxCount         int64
	SecurityGroupIds []string
	SubnetId         string
	Tags             map[*string]*string
}

// NewRunInstancesInput initializes a new RunInstancesInput from the given config.
func (ec *AmazonEC2) NewRunInstancesInput(config RunInstanceConfig) ec2.RunInstancesInput {
	var tags []*ec2.Tag

	for key, v := range config.Tags {
		tags = append(tags, &ec2.Tag{Key: key, Value: v})
	}

	input := ec2.RunInstancesInput{
		DryRun:           aws.Bool(config.DryRun),
		KeyName:          aws.String(config.KeyName),
		MinCount:         aws.Int64(config.MinCount),
		MaxCount:         aws.Int64(config.MaxCount),
		SecurityGroupIds: aws.StringSlice(config.SecurityGroupIds),
		SubnetId:         aws.String(""),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         tags,
			},
		},
	}
	return input
}

// RunInstance requests a single new EC2 instance to AWS.
func (ec *AmazonEC2) RunInstance(ctx context.Context, input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	input.SetDryRun(true)
	for try := 1; try <= ec.Retries; try++ {
		_, err := ec.API.RunInstances(input)
		awsErr, ok := err.(awserr.Error)
		if !ok {
			return nil, err
		}
		if ec.isErrorRetryable(awsErr) {
			logger.Logger(ctx).Info(fmt.Sprintf("[EC2|RUN] Error [%s] when launching nodes on dry mode.\nError: %s\n", awsErr.Code(), awsErr.Message()))
		}
		if ec.isDryRunOperation(awsErr) {
			break
		}
		if try != ec.Retries {
			tools.Sleep(time.Second * time.Duration(try))
			continue
		}
		return nil, err
	}

	input.SetDryRun(false)
	reservation, err := ec.API.RunInstances(input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if !ok {
			return nil, err
		}
		logger.Logger(ctx).Warning(fmt.Sprintf("[EC2|RUN] Error [%s] when launching nodes.\nError: %s\n", awsErr.Code(), awsErr.Message()))
		return nil, err
	}
	return reservation, nil
}

// RunInstances requests a set of new EC2 instances to AWS.
// If there is an error in the middle of the operation, it will return the current reservations as well as the error.
func (ec *AmazonEC2) RunInstances(ctx context.Context, inputs []*ec2.RunInstancesInput) (reservations []*ec2.Reservation, err error) {
	var reservation *ec2.Reservation
	for _, input := range inputs {
		reservation, err = ec.RunInstance(ctx, input)
		if err != nil {
			break
		}
		reservations = append(reservations, reservation)
	}
	return reservations, err
}
