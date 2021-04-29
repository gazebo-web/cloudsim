package factory

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	ErrInvalid = errors.New("invalid")
)

type Test struct {
	Integer int
	String  string
	Bool    bool
	Ptr     *Test
}

func (t *Test) Validate() error {
	if t.Ptr == nil {
		return ErrInvalid
	}

	return nil
}

func TestFactorySuite(t *testing.T) {
	suite.Run(t, new(testFactorySuite))
}

type testFactorySuite struct {
	suite.Suite
	Integer int
	String  string
	Bool    bool
	Ptr     *Test
	m       ConfigValues
	test    Test
}

func (s *testFactorySuite) SetupTest() {
	s.Integer = 123
	s.String = "test"
	s.Bool = true
	s.Ptr = &Test{
		Integer: 456,
		String:  "test_ptr",
		Bool:    false,
	}

	s.m = ConfigValues{
		"integer": s.Integer,
		"string":  s.String,
		"bool":    s.Bool,
		"ptr":     s.Ptr,
	}
	s.test = Test{
		Integer: s.Integer,
		String:  s.String,
		Bool:    s.Bool,
		Ptr:     s.Ptr,
	}
}

func (s *testFactorySuite) TestSetValueMapToStruct() {
	var result Test
	s.Nil(SetValue(&result, s.m))

	// Validate data
	s.Equal(s.test, result)
}

func (s *testFactorySuite) TestSetValueStructToStruct() {
	var result Test
	s.Nil(SetValue(&result, s.test))

	// Validate data
	s.Equal(s.test, result)
}

func (s *testFactorySuite) TestSetValueAndValidateValidStructToStruct() {
	var result Test
	s.Nil(SetValueAndValidate(&result, s.test))

	// Validate data
	s.Equal(s.test, result)
}

func (s *testFactorySuite) TestSetValueAndValidateInvalidStructToStruct() {
	var result Test
	s.test.Ptr = nil
	s.Equal(ErrInvalid, SetValueAndValidate(&result, s.test))
}
