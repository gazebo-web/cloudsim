package mock

import (
	"crypto/rand"
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"strings"
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

	var instances []*ec2.Instance

	for i = 0; i < *input.MaxCount; i++ {
		id := make([]byte, 10)
		if _, err := rand.Read(id); err != nil {
			return nil, err
		}

		var instance ec2.Instance

		for _, tag := range input.TagSpecifications {
			instance.Tags = append(instance.Tags, tag.Tags...)
		}

		instances = append(instances, &instance)

		e.Instances[string(id)] = instance
	}

	return &ec2.Reservation{
		Instances: instances,
	}, nil
}

// TerminateInstances mocks TerminateInstances.
func (e *ec2api) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	for _, id := range input.InstanceIds {
		delete(e.Instances, *id)
	}
	return &ec2.TerminateInstancesOutput{}, nil
}

// DescribeInstances mocks DescribeInstances.
func (e *ec2api) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {

	var found []*ec2.Reservation
	if len(input.Filters) == 0 {
		return &ec2.DescribeInstancesOutput{
			Reservations: nil,
		}, nil
	}

	for _, filter := range input.Filters {
		if !strings.Contains(*filter.Name, "tag:") {
			continue
		}
		tag := strings.Replace(*filter.Name, "tag", "", -1)
		for _, instance := range e.Instances {
			for _, t := range instance.Tags {
				if *t.Key == tag {
					found = append(found, &ec2.Reservation{
						Instances: []*ec2.Instance{
							&instance,
						},
					})
				}
			}
		}
	}

	return &ec2.DescribeInstancesOutput{
		Reservations: found,
	}, nil
}

// WaitUntilInstanceStatusOk mocks WaitUntilInstanceStatusOk.
func (e *ec2api) WaitUntilInstanceStatusOk(input *ec2.DescribeInstanceStatusInput) error {
	return nil
}
