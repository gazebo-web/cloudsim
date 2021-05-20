package mock

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	ec2Impl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines/implementations/ec2"
	"k8s.io/apimachinery/pkg/util/rand"
	"log"
	"strings"
)

// ec2api is an ec2iface.EC2API implementation.
type ec2api struct {
	ec2iface.EC2API
	Instances map[string]ec2.Instance
}

// NewEC2Instance initializes a new ec2.Instance.
func NewEC2Instance(id string, tags []machines.Tag) ec2.Instance {
	var instance ec2.Instance

	instanceID := fmt.Sprintf("i-%s", id)
	instance.InstanceId = &instanceID

	tagSpec := ec2Impl.CreateTagSpecifications(tags)

	var tagList []*ec2.Tag
	for _, tag := range tagSpec {
		tagList = append(tagList, tag.Tags...)
	}

	instance.SetTags(tagList)

	return instance
}

// NewEC2 initializes a new ec2iface.EC2API implementation.
func NewEC2(objects ...ec2.Instance) ec2iface.EC2API {
	instances := make(map[string]ec2.Instance)

	for _, instance := range objects {
		var id string
		_, err := fmt.Sscanf(*instance.InstanceId, "i-%s", &id)
		if err != nil {
			panic(err)
		}
		instances[id] = instance
	}
	return &ec2api{
		Instances: instances,
	}
}

// RunInstances mocks RunInstances.
func (e *ec2api) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if input.DryRun != nil && *input.DryRun {
		return nil, awserr.New("DryRunOperation", "dry run operation", errors.New("dry run error"))
	}

	var i int64

	var instances []*ec2.Instance

	for i = 0; i < *input.MaxCount; i++ {
		var id []rune
		for i := 0; i < 5; i++ {
			id = append(id, rune(rand.Intn(90-65)+65))  // [A-Z]
			id = append(id, rune(rand.Intn(57-48)+48))  // [0-9]
			id = append(id, rune(rand.Intn(122-97)+97)) // [a-z]
		}

		var instance ec2.Instance

		log.Println(id)

		instanceID := fmt.Sprintf("i-%s", string(id))
		instance.InstanceId = &instanceID

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
