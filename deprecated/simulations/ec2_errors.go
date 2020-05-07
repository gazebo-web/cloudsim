package simulations

import "github.com/aws/aws-sdk-go/aws/awserr"

const (
	// AWSErrCodeDryRunOperation is request when a request was a successful dry run that attempted to validate an
	// operation.
	AWSErrCodeDryRunOperation = "DryRunOperation"
	// AWSErrCodeRequestLimitExceeded is returned when too many requests sent to AWS in a short period of time.
	AWSErrCodeRequestLimitExceeded = "RequestLimitExceeded"
	// AWSErrCodeInsufficientInstanceCapacity is returned when not enough instances are available to fulfill the
	// request.
	AWSErrCodeInsufficientInstanceCapacity = "InsufficientInstanceCapacity"
)

var retryableErrors = []string{
	// Dev NOTE: it is important that this array DOES NOT include "DryRunOperation" error.
	AWSErrCodeRequestLimitExceeded,
	AWSErrCodeInsufficientInstanceCapacity,
}

// Deprecated: AWSErrorIsRetryable checks that an error returned by an AWS operation is
// non-fatal and can be retried. This is related to operations requesting
// limited resources on AWS.
func AWSErrorIsRetryable(err awserr.Error) bool {
	for _, awsErrorCode := range retryableErrors {
		if awsErrorCode == err.Code() {
			return true
		}
	}

	return false
}
