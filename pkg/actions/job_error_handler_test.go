package actions

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"testing"
)

var jobErrorTestData = struct {
	// Errors
	fnErr      error
	handlerErr error

	// Job error handlers
	// errHandler correctly handles passed errors and returns no error
	errHandler JobErrorHandler
	// passthroughErrHandler does not handle errors and returns the same error as the function
	passthroughErrHandler JobErrorHandler
	// failingErrHandler does not handle errors and returns a different error from the function
	failingErrHandler JobErrorHandler

	// Job functions
	fn        JobFunc
	failingFn JobFunc

	// Helper functions
	getDeploymentErrorCount func(t *testing.T, db *gorm.DB) int
}{
	// Errors
	fnErr:      errors.New("fn"),
	handlerErr: errors.New("handler"),

	// Job functions
	fn: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
		return value, nil
	},
	failingFn: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
		return value, errors.New("fn")
	},

	// Job error handlers
	errHandler: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
		return value, nil
	},
	passthroughErrHandler: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
		return value, err
	},
	failingErrHandler: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
		return value, errors.New("handler")
	},

	// Helper functions
	getDeploymentErrorCount: func(t *testing.T, db *gorm.DB) int {
		var deploymentCount int
		require.NoError(t, db.Model(&DeploymentError{}).Count(&deploymentCount).Error)

		return deploymentCount
	},
}

func TestJobFuncWrapErrorHandler(t *testing.T) {
	tr := setupTest(t)
	setd := jobErrorTestData

	currentJob := "job_1"
	deployment := &Deployment{CurrentJob: currentJob}

	totalErrCount := 0
	test := func(fn JobFunc, expectedErr error, expectedErrCount int) {
		_, err := fn(tr.ctx, tr.db, deployment, nil)
		if expectedErr != nil {
			require.NotNil(t, err)
			require.Equal(t, expectedErr.Error(), err.Error())
		}

		// Check that the number of job errors has increased as expected
		errCount := setd.getDeploymentErrorCount(t, tr.db)
		require.Equal(t, expectedErrCount, errCount-totalErrCount)
		totalErrCount += expectedErrCount
	}

	// Check that the base functions return the expected values
	test(setd.fn, nil, 0)
	test(setd.failingFn, setd.fnErr, 0)

	// Passing functions
	test(WrapErrorHandler(setd.fn, setd.errHandler), nil, 0)
	test(WrapErrorHandler(setd.fn, setd.passthroughErrHandler), nil, 0)
	test(WrapErrorHandler(setd.fn, setd.failingErrHandler), nil, 0)

	// Failing functions
	test(WrapErrorHandler(setd.failingFn, setd.errHandler), nil, 1)
	test(WrapErrorHandler(setd.failingFn, setd.passthroughErrHandler), setd.fnErr, 1)
	test(WrapErrorHandler(setd.failingFn, setd.failingErrHandler), setd.handlerErr, 2)
}

func TestErrorHandlerIgnoreError(t *testing.T) {
	tr := setupTest(t)
	setd := jobErrorTestData

	currentJob := "job_1"
	deployment := &Deployment{CurrentJob: currentJob}

	test := func(fn JobFunc) {
		wrappedFn := WrapErrorHandler(fn, ErrorHandlerIgnoreError)
		_, err := wrappedFn(tr.ctx, tr.db, deployment, nil)
		require.NoError(t, err)
	}
	test(setd.fn)
	test(setd.failingFn)
}

func TestJobRunErrorHandler(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)
	setd := jobErrorTestData

	deployment, err := newDeployment(tr.db, td.action)
	require.NoError(t, err)

	test := func(job *Job, expectedErr error) {
		_, err := job.Run(tr.ctx, tr.db, deployment, nil)

		// Check error
		if expectedErr != nil {
			require.Error(t, expectedErr, err)
		} else {
			require.NoError(t, err)
		}
	}

	fns := []JobFunc{
		WrapErrorHandler(setd.fn, setd.errHandler),
		WrapErrorHandler(setd.failingFn, setd.errHandler),
	}

	failingFns := append(
		fns,
		WrapErrorHandler(setd.fn, setd.failingErrHandler),
		WrapErrorHandler(setd.failingFn, setd.failingErrHandler),
	)

	// Passing functions
	// PreHooks
	test(&Job{PreHooks: fns, Execute: setd.fn}, nil)
	// Execute
	for _, fn := range fns {
		test(&Job{Execute: fn}, nil)
	}
	// PostHooks
	test(&Job{PostHooks: fns, Execute: setd.fn}, nil)

	// Failing functions
	// PreHooks
	test(&Job{PreHooks: failingFns, Execute: setd.fn}, setd.handlerErr)
	// Execute
	test(&Job{Execute: failingFns[0]}, nil)
	test(&Job{Execute: failingFns[1]}, nil)
	test(&Job{Execute: failingFns[2]}, setd.handlerErr)
	test(&Job{Execute: failingFns[3]}, setd.fnErr)
	// PostHooks
	test(&Job{PostHooks: failingFns, Execute: setd.fn}, setd.handlerErr)
}
