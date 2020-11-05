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
	// AWSErrCodeServiceUnavailable is returned when the request has failed due to a temporary failure of the server.
	AWSErrCodeServiceUnavailable = "ServiceUnavailable"
	// AWSErrCodeUnavailable is returned when the  server is overloaded and can't handle the request.
	AWSErrCodeUnavailable = "Unavailable"
	// AWSErrCodeInternalFailure is returned when the request processing has failed because of an unknown error,
	// exception, or failure.
	AWSErrCodeInternalFailure = "InternalFailure"
	// AWSErrCodeInternalError is returned when an internal error has occurred.
	AWSErrCodeInternalError = "InternalError"
	// AWSErrCodeInsufficientReservedInstanceCapacity is returned when there are not enough available Reserved Instances
	// to satisfy your minimum request.
	AWSErrCodeInsufficientReservedInstanceCapacity = "InsufficientReservedInstanceCapacity"
	// AWSErrCodeInsufficientHostCapacity is returned when there is not enough capacity to fulfill your Dedicated Host request.
	AWSErrCodeInsufficientHostCapacity = "InsufficientHostCapacity"
	// AWSErrCodeInsufficientCapacity is returned when there is not enough capacity to fulfill your import instance request.
	AWSErrCodeInsufficientCapacity = "InsufficientCapacity"
	// AWSErrCodeInsufficientAddressCapacity is returned when not enough available addresses to satisfy your minimum request.
	AWSErrCodeInsufficientAddressCapacity = "InsufficientAddressCapacity"
)

var retryableErrors = []string{
	// Dev NOTE: it is important that this array DOES NOT include "DryRunOperation" error.
	AWSErrCodeRequestLimitExceeded,
	AWSErrCodeInsufficientInstanceCapacity,
	AWSErrCodeServiceUnavailable,
	AWSErrCodeUnavailable,
	AWSErrCodeInternalFailure,
	AWSErrCodeInternalError,
	AWSErrCodeInsufficientReservedInstanceCapacity,
	AWSErrCodeInsufficientHostCapacity,
	AWSErrCodeInsufficientCapacity,
	AWSErrCodeInsufficientAddressCapacity,
}

// AWSErrorIsRetryable checks that an error returned by an AWS operation is
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
