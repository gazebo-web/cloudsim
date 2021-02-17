package defaults

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	defaultValue = "default"
)

// D is a test struct that implements the Defaulter interface.
type D struct {
	Value string
}

// SetDefaults sets default values.
func (d *D) SetDefaults() error {
	d.Value = defaultValue
	return nil
}

// t is a test struct that does not implement the Defaulter interface.
type T struct {
	Value string
}

func TestDefaultsSuite(t *testing.T) {
	suite.Run(t, new(testDefaultsSuite))
}

type testDefaultsSuite struct {
	suite.Suite
}

func (s *testDefaultsSuite) TestSetDefaultsImplements() {
	d := &D{}
	s.Equal("", d.Value)

	s.NoError(SetDefaults(d))

	s.Equal(defaultValue, d.Value)
}

func (s *testDefaultsSuite) TestSetDefaultsDoesNotImplement() {
	t := &T{}
	s.Equal("", t.Value)

	s.NoError(SetDefaults(t))

	s.Equal("", t.Value)
}
