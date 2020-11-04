package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"regexp"
	"time"
)

const (
	// ErrCodeDryRunOperation is the error code returned from a successful dry run operation.
	// A successful dry run operation is returned from AWS when running a request in dry run mode succeeds.
	ErrCodeDryRunOperation = "DryRunOperation"
	// ErrCodeRequestLimitExceeded is returned when too many requests are sent to AWS in a short period of time.
	ErrCodeRequestLimitExceeded = "RequestLimitExceeded"
	// ErrCodeInsufficientInstanceCapacity is returned when not enough instances are available to fulfill the
	// request.
	ErrCodeInsufficientInstanceCapacity = "InsufficientInstanceCapacity"

	// shortSubnetLength specifies the length of a v1 AWS subnet ID.
	shortSubnetLength = 15
	// longSubnetLength specifies the length of a v2 AWS subnet ID.
	longSubnetLength = 24
)

// NewAPI returns an EC2 client from the given config provider.
func NewAPI(config client.ConfigProvider) ec2iface.EC2API {
	return ec2.New(config)
}

// machines is a cloud.Machines implementation.
type machines struct {
	API    ec2iface.EC2API
	Logger ign.Logger
}

// isValidKeyName checks that the given keyName is valid.
func (m *machines) isValidKeyName(keyName string) bool {
	return len(keyName) != 0 && keyName != ""
}

// isValidMachineCount checks that the given min and max values are valid machine count values.
func (m *machines) isValidMachineCount(min, max int64) bool {
	return min > 0 && max > 0 && min <= max
}

// isValidSubnetID checks that the given subnet name is a valid for AWS.
func (m *machines) isValidSubnetID(subnet string) bool {
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
func (m *machines) newRunInstancesInput(createMachines cloud.CreateMachinesInput) *ec2.RunInstancesInput {
	var iamProfile *ec2.IamInstanceProfileSpecification
	iamProfile = &ec2.IamInstanceProfileSpecification{
		Arn:  createMachines.InstanceProfile,
		Name: nil,
	}

	var securityGroups []*string
	for _, sg := range createMachines.FirewallRules {
		securityGroups = append(securityGroups, aws.String(sg))
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
		UserData:          createMachines.InitScript,
		SecurityGroups:    securityGroups,
	}
}

// createTags creates an array of ec2.TagSpecification from the given tag input.
func (m *machines) createTags(input []cloud.Tag) []*ec2.TagSpecification {
	var tagSpec []*ec2.TagSpecification
	for _, tag := range input {
		var tags []*ec2.Tag
		for key, value := range tag.Map {
			tags = append(tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		tagSpec = append(tagSpec, &ec2.TagSpecification{
			ResourceType: aws.String(tag.Resource),
			Tags:         tags,
		})
	}
	return tagSpec
}

// runInstanceDryRun runs a new EC2 instance using dry run mode.
// It will return an cloud.ErrDryRunFailed if the EC2 transaction
// returns a different error code than ErrCodeDryRunOperation.
func (m *machines) runInstanceDryRun(input *ec2.RunInstancesInput) error {
	input.SetDryRun(true)
	_, err := m.API.RunInstances(input)
	awsErr, ok := err.(awserr.Error)
	if !ok || awsErr.Code() != ErrCodeDryRunOperation {
		return cloud.ErrDryRunFailed
	}
	return nil
}

// sleepNSecondsBeforeMaxRetries pauses the current thread for n seconds.
func (m *machines) sleepNSecondsBeforeMaxRetries(n, max int) {
	if n < max {
		time.Sleep(time.Second * time.Duration(n))
	}
}

// parseRunInstanceError parses the given awserr.Error and returns an error from the list of generic errors.
func (m *machines) parseRunInstanceError(err error) error {
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
func (m *machines) runInstance(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	input.SetDryRun(false)
	r, err := m.API.RunInstances(input)
	if err != nil {
		return nil, m.parseRunInstanceError(err)
	}
	return r, nil
}

// create creates a single EC2 instance.
func (m *machines) create(input cloud.CreateMachinesInput) (*cloud.CreateMachinesOutput, error) {
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

	for try := 1; try <= input.Retries; try++ {
		if err := m.runInstanceDryRun(runInstanceInput); err != nil {
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
		return nil, err
	}

	var output cloud.CreateMachinesOutput
	for _, i := range reservation.Instances {
		output.Instances = append(output.Instances, *i.InstanceId)
	}

	return &output, nil
}

// Create creates multiple EC2 instances. It returns the id of the created machines.
// This operation doesn't recover from an error.
// You need to destroy the required machines when an error occurs.
// A single cloud.CreateMachinesOutput instance will be returned for every
// cloud.CreateMachinesInput passed by parameter.
func (m *machines) Create(inputs []cloud.CreateMachinesInput) (created []cloud.CreateMachinesOutput, err error) {
	m.Logger.Debug(fmt.Sprintf("Creating machines with the following input: %+v", inputs))
	var c *cloud.CreateMachinesOutput
	for _, input := range inputs {
		c, err = m.create(input)
		if err != nil {
			m.Logger.Debug(fmt.Sprintf("Creating machines failed while creating the following machine: %+v. Output: %+v. Error: %s", input, created, err))
			return
		}
		created = append(created, *c)
	}
	m.Logger.Debug(fmt.Sprintf("Creating machines succeeded. Output: %+v", created))
	return created, nil
}

// terminate terminates EC2 instances.
// It returns an error if no instances names are provided.
// It also returns an error if the underlying TerminateInstances request fails.
func (m *machines) terminate(input cloud.TerminateMachinesInput) error {
	if input.Instances == nil || len(input.Instances) == 0 {
		return cloud.ErrMissingMachineNames
	}
	terminateInstancesInput := &ec2.TerminateInstancesInput{
		DryRun:      aws.Bool(false),
		InstanceIds: aws.StringSlice(input.Instances),
	}
	_, err := m.API.TerminateInstances(terminateInstancesInput)
	return err
}

// Terminate terminates EC2 machines.
func (m *machines) Terminate(input cloud.TerminateMachinesInput) error {
	m.Logger.Debug(fmt.Sprintf("Terminating machines with the following input: %+v", input))
	err := m.terminate(input)
	if err != nil {
		m.Logger.Debug(fmt.Sprintf("Error while terminating instances. Instances: [%s]. Error: %s.", input.Instances, err))
		return err
	}
	m.Logger.Debug("Terminating machines succeeded.")
	return nil
}

// createFilters creates a set of filters from the given input.
func (m *machines) createFilters(input map[string][]string) []*ec2.Filter {
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
func (m *machines) Count(input cloud.CountMachinesInput) int {
	m.Logger.Debug(fmt.Sprintf("Counting machines with the following parameters: %+v", input))
	filters := m.createFilters(input.Filters)
	out, err := m.API.DescribeInstances(&ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(1000),
		Filters:    filters,
	})
	if err != nil {
		m.Logger.Debug(fmt.Sprintf("Error while counting machines. Error: %s", err))
		return -1
	}
	var count int
	for _, r := range out.Reservations {
		count += len(r.Instances)
	}
	m.Logger.Debug(fmt.Sprintf("Counting machines succedeed. Count: %d", count))
	return count
}

// WaitOK waits for EC2 machines to be in the OK status.
func (m *machines) WaitOK(input []cloud.WaitMachinesOKInput) error {
	m.Logger.Debug(fmt.Sprintf("Waiting for machines to be OK: %+v", input))

	// Collect all instance ids in a single slice
	var instances []string
	for _, i := range input {
		instances = append(instances, i.Instances...)
	}

	// Perform request
	err := m.API.WaitUntilInstanceStatusOk(&ec2.DescribeInstanceStatusInput{
		InstanceIds: aws.StringSlice(instances),
	})
	if err != nil {
		m.Logger.Debug(fmt.Sprintf("Waiting for machines to be OK: %+v failed. Error: %s", input, err))
		return err
	}

	m.Logger.Debug(fmt.Sprintf("Waiting for machines to be OK: %+v succeeded.", input))
	return nil
}

// NewMachines initializes a new cloud.Machines implementation using EC2.
func NewMachines(api ec2iface.EC2API, logger ign.Logger) cloud.Machines {
	return &machines{
		API:    api,
		Logger: logger,
	}
}