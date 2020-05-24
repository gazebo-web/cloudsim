package actions

import (
	"github.com/jinzhu/gorm"
)

// WrapErrorHandler wraps a job function with an ErrorHandler.
// The wrapper also adds any errors returned by the job function or error handler.
// If `fn` returns an error, the error is handled by the `errorHandler` function.
// If the handler returns an error, the error is considered critical and triggers an action execution rollback.
func WrapErrorHandler(fn JobFunc, errorHandler JobErrorHandler) JobFunc {
	return func(tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {

		var err error
		value, err = fn(tx, deployment, value)
		if err != nil {
			// Add the error
			if err := deployment.addJobError(tx, nil, err); err != nil {
				return nil, err
			}

			// Try to handle the error
			var handlerErr error
			value, handlerErr = errorHandler(tx, deployment, value, err)

			// If the handler returned an error, only add it if it differs from the fn error or the same error will be
			// added twice.
			if handlerErr != nil && err.Error() != handlerErr.Error() {
				// Add the error
				if err := deployment.addJobError(tx, nil, err); err != nil {
					return nil, err
				}

			}

			return value, handlerErr
		}

		return value, err
	}
}

// ErrorHandlerIgnoreError ignores errors returned by a function and continues execution.
func ErrorHandlerIgnoreError(tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
	return value, nil
}
