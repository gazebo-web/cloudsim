package errors

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestErrorUtils(t *testing.T) {
	suite.Run(t, new(errorUtilsTestSuite))
}

type errorUtilsTestSuite struct {
	suite.Suite
}

// WARNING: This test is sensitive to the line number it runs in.
// Any new tests must be added below this test.
func (s *errorUtilsTestSuite) TestWithFunctionContext() {
	errMsg := WithFunctionContext(errors.New("test"), "example", 1).Error()
	s.Require().Contains(errMsg, "TestWithFunctionContext")
}
