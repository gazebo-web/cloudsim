package cloud

import "github.com/aws/aws-sdk-go/aws/awserr"

// isErrorRetryable returns true if the code from the given error is retryable.
func (ec *AmazonEC2) isErrorRetryable(err awserr.Error) bool {
	return ec.isCodeRetryable(err.Code())
}

// isDryRunOperation checks that the given error's code is from a Dry run operation.
func (ec *AmazonEC2) isDryRunOperation(err awserr.Error) bool {
	return err.Code() == "DryRunOperation"
}

// isCodeRetryable receives an error code and returns true or false depending on the code's reason.
func (ec *AmazonEC2) isCodeRetryable(code string) bool {
	switch code {
	case "RequestLimitExceeded":
	case "InsufficientInstanceCapacity":
		return true
	}
	return false
}
