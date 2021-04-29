package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
)

// WithFunctionContext wraps an error with information about the function that generated the error.
// `skip` is the number of stack frames that need to be removed to reach the function that generated the error. If this
// function is called to wrap an error on the function that generated the error, then `skip` should be set to 1.
// Higher values of `skip` are necessary to create error wrapping functions that make use of this function.
func WithFunctionContext(err error, errMsg string, skip int) error {
	// Get information about the function calling this function
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return err
	}
	// Get the name of the function that generated the error
	fn := runtime.FuncForPC(pc).Name()
	// Prepare the error message
	msg := fmt.Sprintf("(%s)", fn)
	if errMsg != "" {
		msg = fmt.Sprintf("(%s) %s", fn, errMsg)
	}

	return errors.Wrap(err, msg)
}
