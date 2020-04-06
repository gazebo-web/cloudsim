package cloud

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
)

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
