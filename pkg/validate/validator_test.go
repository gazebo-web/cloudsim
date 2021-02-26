package validate

import (
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

const (
	defaultValue = "default"
)

// V is a test struct that implements the Validator interface.
type V struct {
	Value string `validate:"required"`
}

// Validate validates values.
func (d *V) Validate() error {
	return validator.New().Struct(d)
}

// T is a test struct that does not implement the Validator interface.
type T struct {
	Value string `validate:"required"`
}

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, new(testValidatorSuite))
}

type testValidatorSuite struct {
	suite.Suite
}

func (s *testValidatorSuite) TestValidateImplementsValidData() {
	// The struct implements Validator and contains valid data
	d := &V{
		Value: defaultValue,
	}
	s.Nil(Validate(d))
}

func (s *testValidatorSuite) TestValidateImplementsInvalidData() {
	// The struct implements Validator and contains invalid data
	d := &V{}
	s.NotNil(Validate(d))
}

func (s *testValidatorSuite) TestValidateDoesNotImplement() {
	t := &T{}

	// Struct does not implement Validator, it should not fail the validation
	s.Nil(Validate(t))
}
