package cloud

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
)

func (ec AmazonEC2) CountInstances(ctx context.Context) int {
	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(1000),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:cloudsim-simulation-worker"),
				Values: []*string{
					aws.String(ec.NamePrefix),
				},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("pending"),
					aws.String("running"),
				},
			},
		},
	}

	output, err := ec.API.DescribeInstances(input)
	if err != nil {
		logger.Logger(ctx).Warning("[EC2|COUNT] Error getting the list of available machines.")
		return 0
	}
	return len(output.Reservations)
}