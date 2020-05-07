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

type IAmazonEC2 interface {
	CountInstances(ctx context.Context) int
	TerminateInstances(ctx context.Context, instances []*string) (*ec2.TerminateInstancesOutput, error)
	NewRunInstancesInput(config RunInstancesConfig) ec2.RunInstancesInput
	RunInstance(ctx context.Context, input *ec2.RunInstancesInput) (*ec2.Reservation, error)
	RunInstances(ctx context.Context, inputs []*ec2.RunInstancesInput) (reservations []*ec2.Reservation, err error)
}

// AmazonEC2 wraps the AWS EC2 API.
type AmazonEC2 struct {
	API        ec2iface.EC2API
	Retries    int
	NamePrefix string
}

// NewAmazonEC2 returns a new AmazonEC2 instance by the given AWS session and configuration.
func NewAmazonEC2(p client.ConfigProvider, cfgs ...*aws.Config) IAmazonEC2 {
	var instance AmazonEC2
	if !reflect.ValueOf(p).IsNil() {
		instance.API = ec2.New(p, cfgs...)
	}
	return &instance
}

// CountInstances returns the number of instances that are in both running and pending status.
func (ec *AmazonEC2) CountInstances(ctx context.Context) int {
	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(1000),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:cloudsim-simulation-worker"),
				Values: []*string{
					aws.String(ec.NamePrefix),
				},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("pending"),
					aws.String("running"),
				},
			},
		},
	}

	output, err := ec.API.DescribeInstances(input)
	if err != nil {
		logger.Logger(ctx).Warning("[EC2|COUNT] Error getting the list of available machines.")
		return 0
	}
	return len(output.Reservations)
}

// TerminateInstances terminates a set of EC2 instances by the given instances IDs.
func (ec *AmazonEC2) TerminateInstances(ctx context.Context, instances []*string) (*ec2.TerminateInstancesOutput, error) {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: instances,
	}
	input.SetDryRun(true)
	_, err := ec.API.TerminateInstances(input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if !ok {
			return nil, err
		}
		if !ec.isDryRunOperation(awsErr) {
			logger.Logger(ctx).Info(fmt.Sprintf("[EC2|TERMINATE] Error [%s] while terminating nodes on dry mode.\nError: %s\n", awsErr.Code(), awsErr.Message()))
			return nil, err
		}
	}
	input.SetDryRun(false)
	output, err := ec.API.TerminateInstances(input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if !ok {
			return nil, err
		}
		logger.Logger(ctx).Warning(fmt.Sprintf("[EC2|TERMINATE] Error [%s] while terminating nodes.\nError: %s\n", awsErr.Code(), awsErr.Message()))
		return nil, err
	}
	return output, nil
}

// RunInstancesConfig describes the argument to initialize a new group of EC2 instances.
type RunInstancesConfig struct {
	DryRun           bool
	KeyName          string
	MinCount         int64
	MaxCount         int64
	SecurityGroupIds []string
	SubnetId         string
	Tags             map[*string]*string
}

// NewRunInstancesInput initializes a new RunInstancesInput from the given config.
func (ec *AmazonEC2) NewRunInstancesInput(config RunInstancesConfig) ec2.RunInstancesInput {
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

