package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"regexp"
	"time"
)

const (
	// ErrCodeDryRunOperation is the error code from a successful dry run operation.
	ErrCodeDryRunOperation = "DryRunOperation"
	// ErrCodeRequestLimitExceeded is returned when too many requests sent to AWS in a short period of time.
	ErrCodeRequestLimitExceeded = "RequestLimitExceeded"
	// ErrCodeInsufficientInstanceCapacity is returned when not enough instances are available to fulfill the
	// request.
	ErrCodeInsufficientInstanceCapacity = "InsufficientInstanceCapacity"

	shortSubnetLength = 15
	longSubnetLength  = 24
)

// machines is a cloud.Machines implementation.
type machines struct {
	API ec2iface.EC2API
}

// isValidKeyName checks that the given keyName is valid.
func (m machines) isValidKeyName(keyName string) bool {
	return len(keyName) != 0 && keyName != ""
}

// isValidMachineCount checks that the given min and max values are valid machine count values.
func (m machines) isValidMachineCount(min, max int64) bool {
	return min > 0 && max > 0 && min <= max
}

// isValidSubnetID checks that the given subnet is a valid AWS subnet.
func (m machines) isValidSubnetID(subnet string) bool {
	length := len(subnet)
	if length != shortSubnetLength && length != longSubnetLength {
		return false
	}
	input := []byte(subnet)
	matched, err := regexp.Match("subnet-(\\w+)", input)
	if err != nil {
		return false
	}
	return matched
}

// newRunInstancesInput initializes the configuration to run EC2 instances with the given input.
func (m machines) newRunInstancesInput(createMachines cloud.CreateMachinesInput) *ec2.RunInstancesInput {
	var script *string
	if len(createMachines.InitScript) > 0 {
		script = &createMachines.InitScript
	}

	var iamProfile *ec2.IamInstanceProfileSpecification
	if len(createMachines.ResourceName) > 0 {
		iamProfile = &ec2.IamInstanceProfileSpecification{
			Arn:  &createMachines.ResourceName,
			Name: nil,
		}
	}

	var securitGroups []*string
	for _, sg := range createMachines.FirewallRules {
		securitGroups = append(securitGroups, aws.String(sg))
	}

	tagSpec := m.createTags(createMachines.Tags)
	return &ec2.RunInstancesInput{
		InstanceType:       aws.String(createMachines.Type),
		ImageId:            aws.String(createMachines.Image),
		IamInstanceProfile: iamProfile,
		KeyName:            aws.String(createMachines.KeyName),
		MaxCount:           aws.Int64(createMachines.MaxCount),
		MinCount:           aws.Int64(createMachines.MinCount),
		SecurityGroupIds:   aws.StringSlice(createMachines.FirewallRules),
		SubnetId:           aws.String(createMachines.SubnetID),
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(createMachines.Zone),
		},
		TagSpecifications: tagSpec,
		UserData:          script,
		SecurityGroups:    securitGroups,
	}
}

// createTags creates an array of ec2.TagSpecification from the given tag input.
func (m machines) createTags(input map[string]map[string]string) []*ec2.TagSpecification {
	var tagSpec []*ec2.TagSpecification
	for resource, ts := range input {
		var tags []*ec2.Tag
		for key, value := range ts {
			tags = append(tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		tagSpec = append(tagSpec, &ec2.TagSpecification{
			ResourceType: aws.String(resource),
			Tags:         tags,
		})
	}
	return tagSpec
}

// createInstanceDryRun runs a new EC2 instance using dry run mode.
// It will return an cloud.ErrDryRunFailed if the EC2 transaction
// returns a different error code than ErrCodeDryRunOperation.
func (m machines) createInstanceDryRun(input *ec2.RunInstancesInput) error {
	input.SetDryRun(true)
	_, err := m.API.RunInstances(input)
	awsErr, ok := err.(awserr.Error)
	if !ok || awsErr.Code() != ErrCodeDryRunOperation {
		return cloud.ErrDryRunFailed
	}
	return nil
}

// sleepNSecondsBeforeMaxRetries pauses the current thread for n seconds.
func (m machines) sleepNSecondsBeforeMaxRetries(n, max int) {
	if n < max {
		time.Sleep(time.Second * time.Duration(n))
	}
}

// parseRunInstanceError parses the given awserr.Error and returns an error from the list of generic errors.
func (m machines) parseRunInstanceError(err error) error {
	awsErr, ok := err.(awserr.Error)
	if !ok {
		return cloud.ErrUnknown
	}
	switch awsErr.Code() {
	case ErrCodeInsufficientInstanceCapacity:
		return cloud.ErrInsufficientMachines
	case ErrCodeRequestLimitExceeded:
		return cloud.ErrRequestsLimitExceeded
	}
	return cloud.ErrUnknown
}

// runInstance is a wrapper for the EC2 RunInstances method.
func (m machines) runInstance(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	input.SetDryRun(false)
	r, err := m.API.RunInstances(input)
	if err != nil {
		return nil, m.parseRunInstanceError(err)
	}
	return r, nil
}

// create creates a single EC2 instance.
func (m machines) create(input cloud.CreateMachinesInput) (*cloud.CreateMachinesOutput, error) {
	if !m.isValidKeyName(input.KeyName) {
		return nil, cloud.ErrMissingKeyName
	}
	if !m.isValidMachineCount(input.MinCount, input.MaxCount) {
		return nil, cloud.ErrInvalidMachinesCount
	}
	if !m.isValidSubnetID(input.SubnetID) {
		return nil, cloud.ErrInvalidSubnetID
	}

	runInstanceInput := m.newRunInstancesInput(input)

	var output cloud.CreateMachinesOutput
	for try := 1; try <= input.Retries; try++ {
		if err := m.createInstanceDryRun(runInstanceInput); err != nil {
			m.sleepNSecondsBeforeMaxRetries(try, input.Retries)
		} else {
			break
		}
		if try == input.Retries {
			return nil, cloud.ErrUnknown
		}
	}

	reservation, err := m.runInstance(runInstanceInput)
	if err != nil {
		return &output, err
	}

	for _, i := range reservation.Instances {
		output.Instances = append(output.Instances, *i.InstanceId)
	}

	return &output, nil
}

// Create creates multiple EC2 instances. It will return the created machines.
// This operation doesn't recover from an error.
// You need to destroy the required machines when an error occurs.
func (m machines) Create(inputs []cloud.CreateMachinesInput) (created []cloud.CreateMachinesOutput, err error) {
	var c *cloud.CreateMachinesOutput
	for _, input := range inputs {
		c, err = m.create(input)
		if err != nil {
			return
		}
		created = append(created, *c)
	}
	return created, nil
}

// terminateInstanceDryRun terminates an instance using dry run mode.
func (m machines) terminateInstanceDryRun(input *ec2.TerminateInstancesInput) error {
	input.SetDryRun(true)
	_, err := m.API.TerminateInstances(input)
	awsErr, ok := err.(awserr.Error)
	if !ok || awsErr.Code() != ErrCodeDryRunOperation {
		return cloud.ErrDryRunFailed
	}
	return nil
}

// Terminate terminates EC2 machines.
func (m machines) Terminate(input cloud.TerminateMachinesInput) error {
	if input.Names == nil || len(input.Names) == 0 {
		return cloud.ErrMissingMachineNames
	}

	terminateInstancesInput := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice(input.Names),
	}

	for try := 1; try <= input.Retries; try++ {
		if err := m.terminateInstanceDryRun(terminateInstancesInput); err != nil {
			m.sleepNSecondsBeforeMaxRetries(try, input.Retries)
		} else {
			break
		}
		if try == input.Retries {
			return cloud.ErrUnknown
		}
	}

	terminateInstancesInput.SetDryRun(false)
	_, err := m.API.TerminateInstances(terminateInstancesInput)
	if err != nil {
		return err
	}
	return nil
}

func (m machines) createFilters(input map[string][]string) []*ec2.Filter {
	var filters []*ec2.Filter
	for k, v := range input {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(k),
			Values: aws.StringSlice(v),
		})
	}
	return filters
}

// Count counts EC2 machines.
func (m machines) Count(input cloud.CountMachinesInput) int {
	if input.MaxResults > 1000 || input.MaxResults < 5 {
		return -1
	}
	filters := m.createFilters(input.Filters)
	out, err := m.API.DescribeInstances(&ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(int64(input.MaxResults)),
		Filters:    filters,
	})
	if err != nil {
		return -1
	}
	var count int
	for _, r := range out.Reservations {
		count += len(r.Instances)
	}
	return count
}

// NewMachines initializes a new cloud.Machines implementation using EC2.
func NewMachines(api ec2iface.EC2API) cloud.Machines {
	return &machines{
		API: api,
	}
}
