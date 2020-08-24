package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/simulations"
)

const (
	Ec2OpRunInstances              OpType = "RunInstances"
	Ec2OpStopInstances             OpType = "StopInstances"
	Ec2OpTerminateInstances        OpType = "TerminateInstances"
	Ec2OpWaitUntilInstanceStatusOk OpType = "WaitUntilInstanceStatusOk"
)

// EC2Mock is a mock for EC2 service
type EC2Mock struct {
	ec2iface.EC2API
	Mock
	RunInstancesFunc              func(*ec2.RunInstancesInput) (*ec2.Reservation, error)
	StopInstancesFunc             func(*ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error)
	TerminateInstancesFunc        func(*ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error)
	WaitUntilInstanceStatusOkFunc func(*ec2.DescribeInstanceStatusInput) error
}

// NewEC2Mock creates a new EC2Mock.
func NewEC2Mock() *EC2Mock {
	m := &EC2Mock{}
	m.Reset()
	return m
}

// NewEC2MockSuccessfulLaunch creates a new Mock prepared to launch EC2 instances required for a single robot
// simulation.
func NewEC2MockSuccessfulLaunch() *EC2Mock {
	m := NewEC2Mock()
	m.SetMockFunction(Ec2OpWaitUntilInstanceStatusOk, FixedValues, false, nil)
	m.SetMockFunction(Ec2OpRunInstances, FixedValues, false,
		// Check for available machines
		m.NewAWSErr(simulations.AWSErrCodeDryRunOperation),
		// EC2 Instance 1
		m.NewAWSErr(simulations.AWSErrCodeDryRunOperation),
		m.NewReservation(fmt.Sprintf("i-test-1-%s", uuid.NewV4().String())),
		// EC2 Instance 2
		m.NewAWSErr(simulations.AWSErrCodeDryRunOperation),
		m.NewReservation(fmt.Sprintf("i-test-2-%s", uuid.NewV4().String())),
	)
	return m
}

// NewAWSErr is a helper method to create AWS errors.
func (m *EC2Mock) NewAWSErr(code string) awserr.Error {
	return awserr.New(code, code, nil)
}

// NewInstance is a helper method to easily create ec2.Instance structs.
func (m *EC2Mock) NewInstance(iID string) *ec2.Instance {
	return &ec2.Instance{
		InstanceId: &iID,
	}
}

// NewReservation is a helper method to easily create "RunInstances" results.
func (m *EC2Mock) NewReservation(iID ...string) *ec2.Reservation {
	instances := make([]*ec2.Instance, len(iID))
	for i, id := range iID {
		instances[i] = m.NewInstance(id)
	}
	return &ec2.Reservation{
		Instances: instances,
	}
}

// RunInstances API operation for Amazon Elastic Compute Cloud.
func (m *EC2Mock) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if m.RunInstancesFunc != nil {
		return m.RunInstancesFunc(input)
	}

	defer m.InvokeCallback(Ec2OpRunInstances, input)

	result := m.GetMockResult(Ec2OpRunInstances)
	// PassThrough is a special value that indicates the non-mocked version of
	// this function should be called
	if result == PassThrough {
		return m.EC2API.RunInstances(input)
	}
	// If the mock result is an error, return that error
	if err, ok := result.(error); ok {
		return nil, err
	}

	// This explicit cast is needed to avoid a panic when result is 'nil'.
	if r, ok := result.(*ec2.Reservation); ok {
		return r, nil
	}
	return nil, nil
}

// StopInstances API operation for Amazon Elastic Compute Cloud.
func (m *EC2Mock) StopInstances(input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	if m.StopInstancesFunc != nil {
		return m.StopInstancesFunc(input)
	}

	defer m.InvokeCallback(Ec2OpStopInstances, input)

	result := m.GetMockResult(Ec2OpStopInstances)
	// PassThrough is a special value that indicates the non-mocked version of
	// this function should be called
	if result == PassThrough {
		return m.EC2API.StopInstances(input)
	}
	// If the mock result is an error, return that error
	if err, ok := result.(error); ok {
		return nil, err
	}

	if r, ok := result.(*ec2.StopInstancesOutput); ok {
		return r, nil
	}
	return nil, nil
}

// TerminateInstances API operation for Amazon Elastic Compute Cloud.
func (m *EC2Mock) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	if m.TerminateInstancesFunc != nil {
		return m.TerminateInstancesFunc(input)
	}

	defer m.InvokeCallback(Ec2OpTerminateInstances, input)

	result := m.GetMockResult(Ec2OpTerminateInstances)
	// PassThrough is a special value that indicates the non-mocked version of
	// this function should be called
	if result == PassThrough {
		return m.EC2API.TerminateInstances(input)
	}
	// If the mock result is an error, return that error
	if err, ok := result.(error); ok {
		return nil, err
	}

	if r, ok := result.(*ec2.TerminateInstancesOutput); ok {
		return r, nil
	}
	return nil, nil
}

// WaitUntilInstanceStatusOk uses the Amazon EC2 API operation
// DescribeInstanceStatus to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (m *EC2Mock) WaitUntilInstanceStatusOk(input *ec2.DescribeInstanceStatusInput) error {
	if m.WaitUntilInstanceStatusOkFunc != nil {
		return m.WaitUntilInstanceStatusOkFunc(input)
	}

	defer m.InvokeCallback(Ec2OpWaitUntilInstanceStatusOk, input)

	result := m.GetMockResult(Ec2OpWaitUntilInstanceStatusOk)
	// PassThrough is a special value that indicates the non-mocked version of
	// this function should be called
	if result == PassThrough {
		return m.EC2API.WaitUntilInstanceStatusOk(input)
	}
	// If the mock result is an error, return that error
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

// ModifyInstanceAttribute mocks the EC2 API operation to return successfully.
// Mocking return values for this method is not currently supported.
// Modifies an EC2 instance attribute.
func (m *EC2Mock) ModifyInstanceAttribute(input *ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput,
	error) {
	return nil, nil
}
