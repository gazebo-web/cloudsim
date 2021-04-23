package ec2

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"regexp"
	"strings"
	"sync"
	"text/template"
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

// ec2Machines is a machines.Machines implementation.
type ec2Machines struct {
	API             ec2iface.EC2API
	Logger          ign.Logger
	workerGroupName string
	limit           int64
	lock            sync.Mutex
}

// isValidKeyName checks that the given keyName is valid.
func (m *ec2Machines) isValidKeyName(keyName string) bool {
	return len(keyName) != 0 && keyName != ""
}

// isValidMachineCount checks that the given min and max values are valid machine count values.
func (m *ec2Machines) isValidMachineCount(min, max int64) bool {
	return min > 0 && max > 0 && min <= max
}

// isValidSubnetID checks that the given subnet name is a valid for AWS.
func (m *ec2Machines) isValidSubnetID(subnet string) bool {
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

// isValidClusterID checks that the given cluster ID is valid.
func (m *ec2Machines) isValidClusterID(clusterID string) bool {
	return len(clusterID) > 0
}

// newRunInstancesInput initializes the configuration to run EC2 instances with the given input.
func (m *ec2Machines) newRunInstancesInput(createMachines machines.CreateMachinesInput) *ec2.RunInstancesInput {
	var iamProfile *ec2.IamInstanceProfileSpecification
	iamProfile = &ec2.IamInstanceProfileSpecification{
		Arn:  createMachines.InstanceProfile,
		Name: nil,
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
	}
}

// createTags creates an array of ec2.TagSpecification from the given tag input.
func (m *ec2Machines) createTags(input []machines.Tag) []*ec2.TagSpecification {
	return CreateTagSpecifications(input)
}

// CreateTagSpecifications converts a set of tags into ec2 tag specifications.
func CreateTagSpecifications(input []machines.Tag) []*ec2.TagSpecification {
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
// It will return an machines.ErrDryRunFailed if the EC2 transaction
// returns a different error code than ErrCodeDryRunOperation.
func (m *ec2Machines) runInstanceDryRun(input *ec2.RunInstancesInput) error {
	input.SetDryRun(true)
	_, err := m.API.RunInstances(input)
	awsErr, ok := err.(awserr.Error)
	if !ok || awsErr.Code() != ErrCodeDryRunOperation {
		return errors.Wrap(machines.ErrUnknown, err.Error())
	}
	return nil
}

// sleepNSecondsBeforeMaxRetries pauses the current thread for n seconds.
func (m *ec2Machines) sleepNSecondsBeforeMaxRetries(n, max int) {
	if n < max {
		time.Sleep(time.Second * time.Duration(n))
	}
}

// parseRunInstanceError parses the given awserr.Error and returns an error from the list of generic errors.
func (m *ec2Machines) parseRunInstanceError(err error) error {
	awsErr, ok := err.(awserr.Error)
	if !ok {
		return machines.ErrUnknown
	}
	switch awsErr.Code() {
	case ErrCodeInsufficientInstanceCapacity:
		return machines.ErrInsufficientMachines
	case ErrCodeRequestLimitExceeded:
		return machines.ErrRequestsLimitExceeded
	}
	return errors.Wrap(machines.ErrUnknown, err.Error())
}

// runInstance is a wrapper for the EC2 RunInstances method.
func (m *ec2Machines) runInstance(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	input.SetDryRun(false)
	r, err := m.API.RunInstances(input)
	if err != nil {
		return nil, m.parseRunInstanceError(err)
	}
	return r, nil
}

// create creates a single EC2 instance.
func (m *ec2Machines) create(input machines.CreateMachinesInput) (*machines.CreateMachinesOutput, error) {
	if !m.isValidKeyName(input.KeyName) {
		return nil, machines.ErrMissingKeyName
	}
	if !m.isValidMachineCount(input.MinCount, input.MaxCount) {
		return nil, machines.ErrInvalidMachinesCount
	}
	if !m.isValidSubnetID(input.SubnetID) {
		return nil, machines.ErrInvalidSubnetID
	}
	if !m.isValidClusterID(input.ClusterID) {
		return nil, machines.ErrInvalidClusterID
	}
	if !m.checkAvailableMachines(input.MaxCount) {
		return nil, machines.ErrInsufficientMachines
	}

	if input.InitScript == nil {
		userData, err := m.createUserData(input)
		if err != nil {
			return nil, err
		}
		// EC2 requires that user data strings are encoded in base64
		userData = base64.StdEncoding.EncodeToString([]byte(userData))
		input.InitScript = &userData
	}

	runInstanceInput := m.newRunInstancesInput(input)

	for try := 1; try <= input.Retries; try++ {
		var err error
		if err = m.runInstanceDryRun(runInstanceInput); err != nil {
			m.sleepNSecondsBeforeMaxRetries(try, input.Retries)
		} else {
			break
		}
		if try == input.Retries {
			return nil, errors.Wrap(machines.ErrUnknown, fmt.Sprintf("max retries, with error: %s", err.Error()))
		}
	}

	reservation, err := m.runInstance(runInstanceInput)
	if err != nil {
		return nil, err
	}

	var output machines.CreateMachinesOutput
	for _, i := range reservation.Instances {
		output.Instances = append(output.Instances, *i.InstanceId)
	}

	return &output, nil
}

// Create creates multiple EC2 instances. It returns the id of the created machines.
// This operation doesn't recover from an error.
// You need to destroy the required machines when an error occurs.
// A single machines.CreateMachinesOutput instance will be returned for every
// machines.CreateMachinesInput passed by parameter.
func (m *ec2Machines) Create(inputs []machines.CreateMachinesInput) (created []machines.CreateMachinesOutput, err error) {
	m.Logger.Debug(fmt.Sprintf("Creating machines with the following input: %+v", inputs))

	// A lock is used to synchronize multiple worker requests.
	m.lock.Lock()
	defer m.lock.Unlock()

	var c *machines.CreateMachinesOutput
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

// terminateByID terminates EC2 instances by instance ID.
// It returns an error if no instance ids are provided.
// It also returns an error if the underlying TerminateInstances request fails.
func (m *ec2Machines) terminateByID(instances []string) error {
	_, err := m.API.TerminateInstances(&ec2.TerminateInstancesInput{
		DryRun:      aws.Bool(false),
		InstanceIds: aws.StringSlice(instances),
	})
	return err
}

// terminateByFilters terminates EC2 instances by filtering.
// It returns an error if no instances filters are provided.
// It also returns an error if the underlying TerminateInstances request fails.
func (m *ec2Machines) terminateByFilters(filters map[string][]string) error {
	out, err := m.API.DescribeInstances(&ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(1000),
		Filters:    m.createFilters(filters),
	})
	if err != nil {
		return err
	}

	if len(out.Reservations) == 0 {
		return nil
	}

	var instanceIds []string
	for _, r := range out.Reservations {
		for _, instance := range r.Instances {
			instanceIds = append(instanceIds, *instance.InstanceId)
		}
	}

	return m.terminateByID(instanceIds)
}

// Terminate terminates EC2 machines by either passing instances ids or filters.
// If both are passed, instances will be terminated using instances ids first, and then filters.
// If the former fails, the latter won't be executed.
func (m *ec2Machines) Terminate(input machines.TerminateMachinesInput) error {
	m.Logger.Debug(fmt.Sprintf("Terminating machines with the following input: %+v", input))

	err := input.Validate()
	if err != nil {
		m.Logger.Debug(fmt.Sprintf("Invalid request, couldn't validate input: %+v", input))
		return err
	}

	if input.ValidateInstances() == nil {
		err = m.terminateByID(input.Instances)
		if err != nil {
			m.Logger.Debug(fmt.Sprintf("Error while terminating instances by id. IDs: [%s]. Error: %s.", input.Filters, err))
			return err
		}
	}

	if input.ValidateFilters() == nil {
		err = m.terminateByFilters(input.Filters)
		if err != nil {
			m.Logger.Debug(fmt.Sprintf("Error while terminating instances by filters. Filters: [%s]. Error: %s.", input.Filters, err))
			return err
		}
	}

	m.Logger.Debug("Terminating machines succeeded.")
	return nil
}

// createFilters creates a set of filters from the given input.
func (m *ec2Machines) createFilters(input map[string][]string) []*ec2.Filter {
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
func (m *ec2Machines) Count(input machines.CountMachinesInput) int {
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
func (m *ec2Machines) WaitOK(input []machines.WaitMachinesOKInput) error {
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

// createUserData generates a bash command to make the node join the cluster.
func (m *ec2Machines) createUserData(input machines.CreateMachinesInput) (string, error) {
	tmpl, err := template.ParseFiles("ec2_user_data.sh")
	if err != nil {
		return "", err
	}

	var b []byte
	buffer := bytes.NewBuffer(b)

	labels := make([]string, 0, len(input.Labels))

	for k, v := range input.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}

	err = tmpl.Execute(buffer, map[string]interface{}{
		"Labels":      strings.Join(labels, ","),
		"ClusterName": input.ClusterID,
		"Args":        "--use-max-pods false",
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (m *ec2Machines) checkAvailableMachines(requested int64) bool {
	// If limit is set to a number lower than zero, it means that there is no limit for machines.
	if m.limit < 0 {
		return true
	}

	// Get the number of provisioned machines from cloud provider.
	count := m.Count(machines.CountMachinesInput{
		Filters: map[string][]string{
			"tag:cloudsim-simulation-worker": {
				m.workerGroupName,
			},
			"instance-state-name": {
				"pending",
				"running",
			},
		},
	})

	return requested <= m.limit-int64(count)
}

// List is used to list all pending, running, shutting-down, stopping, stopped and terminated instances with their respective status.
func (m *ec2Machines) List(input machines.ListMachinesInput) (*machines.ListMachinesOutput, error) {
	m.Logger.Debug(fmt.Sprintf("Listing machines with the following input: %+v", input))
	res, err := m.API.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		Filters:             m.createFilters(input.Filters),
		IncludeAllInstances: aws.Bool(true),
		MaxResults:          aws.Int64(1000),
	})
	if err != nil {
		m.Logger.Debug(fmt.Sprintf("Listing machines with the following input: %+v failed, error: %s", input, err))
		return nil, err
	}

	var output machines.ListMachinesOutput
	output.Instances = make([]machines.ListMachinesItem, len(res.InstanceStatuses))

	for i, instanceStatus := range res.InstanceStatuses {
		output.Instances[i] = machines.ListMachinesItem{
			InstanceID: *instanceStatus.InstanceId,
			State:      *instanceStatus.InstanceState.Name,
		}
	}

	m.Logger.Debug(fmt.Sprintf("Listing machines with the following input: %+v succeded. Output: %+v", input, output))

	return &output, nil
}

// MachinesConfig includes a set of field to configure a machines.Machines implementation using EC2.
type MachinesConfig struct {
	// API has a reference to the EC2 API.
	API ec2iface.EC2API
	// Logger is an instance of ign.Logger for logging messages in the Machines component.
	Logger ign.Logger
	// Limit sets the limit of machines that can be created, if this value is nil, no limit will be set.
	Limit *int64
	// WorkerGroupName is label set to all machines created by the Machines component and it's used when counting the
	// amount of machines available.
	WorkerGroupName string
}

// NewMachinesWithConfig initializes a new machines.Machines implementation configured by the given MachinesConfig.
func NewMachinesWithConfig(cfg MachinesConfig) machines.Machines {
	limit := int64(-1)
	if cfg.Limit != nil && *cfg.Limit >= 0 {
		limit = *cfg.Limit
	}

	return &ec2Machines{
		API:             cfg.API,
		Logger:          cfg.Logger,
		limit:           limit,
		workerGroupName: cfg.WorkerGroupName,
	}
}

// NewMachines initializes a new machines.Machines implementation using EC2 with some default configuration:
//	* Limit: -1
//	* Worker group name: "cloudsim-simulation-worker"
func NewMachines(api ec2iface.EC2API, logger ign.Logger) machines.Machines {
	return &ec2Machines{
		API:             api,
		Logger:          logger,
		limit:           -1,
		workerGroupName: "cloudsim-simulation-worker",
	}
}
