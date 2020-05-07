package cloud

import (
	"context"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/mock"
)

type AmazonEC2Mock struct {
	*mock.Mock
}

func NewAmazonEC2Mock() *AmazonEC2Mock {
	var ec2 *AmazonEC2Mock
	ec2 = &AmazonEC2Mock{
		Mock: new(mock.Mock),
	}
	return ec2
}

func (ec *AmazonEC2Mock) CountInstances(ctx context.Context) int {
	args := ec.Called(ctx)
	return args.Int(0)
}

func (ec *AmazonEC2Mock) TerminateInstances(ctx context.Context, instances []*string) (*ec2.TerminateInstancesOutput, error) {
	args := ec.Called(ctx, instances)
	output := args.Get(0).(*ec2.TerminateInstancesOutput)
	return output, args.Error(1)
}

func (ec *AmazonEC2Mock) NewRunInstancesInput(config RunInstancesConfig) ec2.RunInstancesInput {
	args := ec.Called(config)
	output := args.Get(0).(ec2.RunInstancesInput)
	return output
}

func (ec *AmazonEC2Mock) RunInstance(ctx context.Context, input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	args := ec.Called(ctx, input)
	output := args.Get(0).(*ec2.Reservation)
	return output, args.Error(1)
}

func (ec *AmazonEC2Mock) RunInstances(ctx context.Context, inputs []*ec2.RunInstancesInput) (reservations []*ec2.Reservation, err error) {
	args := ec.Called(ctx, inputs)
	reservations = args.Get(0).([]*ec2.Reservation)
	err = args.Error(1)
	return
}


