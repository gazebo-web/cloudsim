package yaml

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"io/ioutil"
	"os"
	"testing"
)

type TestYAMLObject struct {
	Int   int
	Str   string
	Bool  bool
	Array []int
}

// TestYAML is the target struct to parse the test YAML file into. It tests all of the following:
// * Data types
// * Custom variable names
// * Comments
// * References (& and *)
type TestYAML struct {
	Integer int    `yaml:"int"`
	String  string `yaml:"str"`
	Boolean bool   `yaml:"bool"`
	Array   []int
	Object  TestYAMLObject
}

const (
	testYAML = `
int: &int 1247129
str: &str "test"  
bool: &bool true   
array: &array       
  - 1
  - 2
  - 3
object:
  int: *int
  str: *str
  bool: *bool
  array: *array
# There is no entry in the target struct for this variable
extra: string
`
)

func TestYAMLLoader(t *testing.T) {
	suite.Run(t, new(testLoaderSuite))
}

type testLoaderSuite struct {
	suite.Suite
	loader *loader
}

func (s *testLoaderSuite) SetupSuite() {
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
	s.loader = NewLoader(logger).(*loader)
}

func (s *testLoaderSuite) TestLog() {
	s.NotPanics(func() { s.loader.log("test") })
}

func (s *testLoaderSuite) TestNilLog() {
	// Verify that calling `loader.log` with a `nil` logger does not panic
	loader := NewLoader(nil).(*loader)
	s.NotPanics(func() { loader.log("test") })
}

func (s *testLoaderSuite) TestParse() {
	expected := &TestYAML{
		Integer: 1247129,
		String:  "test",
		Boolean: true,
		Array:   []int{1, 2, 3},
		Object: TestYAMLObject{
			Int:   1247129,
			Str:   "test",
			Bool:  true,
			Array: []int{1, 2, 3},
		},
	}

	out := &TestYAML{}

	err := s.loader.parseBytes([]byte(testYAML), out)
	s.NoError(err)
	s.Equal(expected, out)
}

func (s *testLoaderSuite) TestFailParse() {
	invalidYAML := `a:-`

	out := &TestYAML{}

	err := s.loader.parseBytes([]byte(invalidYAML), out)
	s.Error(err)
}

func (s *testLoaderSuite) TestIncompatibleOutputStruct() {
	type IncompatibleOut struct {
		Test string
	}

	expected := &IncompatibleOut{}

	out := &IncompatibleOut{}

	// Loading a YAML file into an incompatible structure should result in no fields being loaded, but no error
	err := s.loader.parseBytes([]byte(testYAML), out)
	s.NoError(err)
	s.Equal(expected, out)
}

func (s *testLoaderSuite) TestLoad() {
	// Create a YAML file to read from
	filename := "test.yaml"
	if err := ioutil.WriteFile(filename, []byte(testYAML), os.FileMode(0700)); err != nil {
		s.Fail("Failed to write YAML test file.", err)
	}
	// Delete the file after finishing with the test
	defer func() {
		_ = os.Remove(filename)
	}()

	expected := &TestYAML{
		Integer: 1247129,
		String:  "test",
		Boolean: true,
		Array:   []int{1, 2, 3},
		Object: TestYAMLObject{
			Int:   1247129,
			Str:   "test",
			Bool:  true,
			Array: []int{1, 2, 3},
		},
	}

	out := &TestYAML{}

	// Loading a YAML file into an incompatible structure should result in no error and no fields being loaded
	err := s.loader.Load(filename, out)
	s.NoError(err)
	s.Equal(expected, out)
}
