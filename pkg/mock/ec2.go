package mock

import (
	"crypto/rand"
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// ec2api is an ec2iface.EC2API implementation.
type ec2api struct {
	ec2iface.EC2API
	Instances map[string]ec2.Instance
}

// NewEC2 initializes a new ec2iface.EC2API implementation.
func NewEC2() ec2iface.EC2API {
	return &ec2api{
		Instances: make(map[string]ec2.Instance),
	}
}

// RunInstances mocks RunInstances.
func (e *ec2api) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if input.DryRun != nil && *input.DryRun {
		return nil, awserr.New("ErrCodeDryRunOperation", "dry run operation", errors.New("dry run error"))
	}

	var i int64
	for i = 0; i < *input.MaxCount; i++ {
		id := make([]byte, 10)
		if _, err := rand.Read(id); err != nil {
			return nil, err
		}
		e.Instances[string(id)] = ec2.Instance{}
	}
	return &ec2.Reservation{}, nil
}

// TerminateInstances mocks TerminateInstances.
func (e *ec2api) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return &ec2.TerminateInstancesOutput{}, nil
}

// DescribeInstances mocks DescribeInstances.
func (e *ec2api) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return nil, nil
}

// WaitUntilInstanceStatusOk mocks WaitUntilInstanceStatusOk.
func (e *ec2api) WaitUntilInstanceStatusOk(*ec2.DescribeInstanceStatusInput) error {
	return nil
}
