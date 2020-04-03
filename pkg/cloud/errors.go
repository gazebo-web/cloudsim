package cloud

import "github.com/aws/aws-sdk-go/aws/awserr"

func (ec *AmazonEC2) isErrorRetryable(err awserr.Error) bool {
	return ec.isCodeRetryable(err.Code())
}

func (ec *AmazonEC2) isDryRunOperation(err awserr.Error) bool {
	return err.Code() == "DryRunOperation"
}

func (ec *AmazonEC2) isCodeRetryable(code string) bool {
	switch code {
	case "RequestLimitExceeded":
	case "InsufficientInstanceCapacity":
		return true
	}
	return false
}
