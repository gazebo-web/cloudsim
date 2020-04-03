package cloud

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"time"
)

func (ec AmazonEC2) Launch(ctx context.Context, input *ec2.RunInstancesInput) (*ec2.Reservation, error) {

	input.SetDryRun(true)
	for try := 1; try <= ec.Retries; try++ {
		ec.lockLaunch.Lock()
		_, err := ec.API.RunInstances(input)
		awsErr, ok := err.(awserr.Error)
		if !ok {
			ec.lockLaunch.Unlock()
			return nil, err
		}
		if ec.isErrorRetryable(awsErr) {
			logger.Logger(ctx).Info(fmt.Sprintf("[EC2|LAUNCH] Error [%s] while launching nodes on dry mode.\nError: %s\n", awsErr.Code(), awsErr.Message()))
		}
		if ec.isDryRunOperation(awsErr) {
			break
		}
		if try != ec.Retries {
			ec.lockLaunch.Unlock()
			tools.Sleep(time.Second * time.Duration(try))
		}
	}

	input.SetDryRun(false)
	reservation, err := ec.API.RunInstances(input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if !ok {
			return nil, err
		}
		logger.Logger(ctx).Warning(fmt.Sprintf("[EC2|LAUNCH] Error [%s] while launching nodes.\nError: %s\n", awsErr.Code(), awsErr.Message()))
		return nil, err
	}
	return reservation, nil
}