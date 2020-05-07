package simulations

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"reflect"
)

// Deprecated: AWSFactory is the single place where all AWS service instances are created.
type AWSFactory struct {
	isGoTest bool
}

// Deprecated: NewAWSFactory creates a new AWSFactory
func NewAWSFactory(isGoTest bool) *AWSFactory {
	f := AWSFactory{}
	f.isGoTest = isGoTest
	return &f
}

// Deprecated: MockableS3 is a type used in tests to allow for easy mocking of S3 service.
type MockableS3 struct {
	s3iface.S3API
}

// Deprecated: NewS3Svc creates a new instance of the S3 client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
func (f *AWSFactory) NewS3Svc(p client.ConfigProvider, cfgs ...*aws.Config) s3iface.S3API {
	var svc s3iface.S3API
	if !reflect.ValueOf(p).IsNil() {
		svc = s3.New(p, cfgs...)
	}

	if f.isGoTest {
		return &MockableS3{S3API: svc}
	}
	return svc
}

// Deprecated: EnsureMockableS3 casts the given arg to MockableS3 or fails.
func EnsureMockableS3(svc s3iface.S3API) *MockableS3 {
	return svc.(*MockableS3)
}

// Deprecated: MockableEC2 is a type used in tests to allow for easy mocking of EC2 operations.
type MockableEC2 struct {
	ec2iface.EC2API
}

// Deprecated: NewEC2Svc creates a new instance of the EC2 client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
func (f *AWSFactory) NewEC2Svc(p client.ConfigProvider, cfgs ...*aws.Config) ec2iface.EC2API {
	var svc ec2iface.EC2API
	if !reflect.ValueOf(p).IsNil() {
		svc = ec2.New(p, cfgs...)
	}

	if f.isGoTest {
		return &MockableEC2{EC2API: svc}
	}
	return svc
}

// Deprecated: AssertMockedEC2 casts the given arg to MockableEC2 or fails.
func AssertMockedEC2(svc ec2iface.EC2API) *MockableEC2 {
	return svc.(*MockableEC2)
}

// Deprecated: SetImpl sets the underlying implementation of this MockableEC2
func (m *MockableEC2) SetImpl(svc ec2iface.EC2API) {
	m.EC2API = svc
}
