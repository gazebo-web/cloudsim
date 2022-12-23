package loader

import (
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
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
	suite.Run(t, new(testYAMLLoaderSuite))
}

type testYAMLLoaderSuite struct {
	suite.Suite
	loader *yamlLoader
}

func (s *testYAMLLoaderSuite) SetupSuite() {
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
	s.loader = NewYAMLLoader(logger).(*yamlLoader)
}

func (s *testYAMLLoaderSuite) TestLog() {
	s.NotPanics(func() { s.loader.log("test") })
}

func (s *testYAMLLoaderSuite) TestNilLog() {
	// Verify that calling `yamlLoader.log` with a `nil` logger does not panic
	loader := NewYAMLLoader(nil).(*yamlLoader)
	s.NotPanics(func() { loader.log("test") })
}

func (s *testYAMLLoaderSuite) TestParse() {
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

func (s *testYAMLLoaderSuite) TestFailParse() {
	invalidYAML := `a:-`

	out := &TestYAML{}

	err := s.loader.parseBytes([]byte(invalidYAML), out)
	s.Error(err)
}

func (s *testYAMLLoaderSuite) TestIncompatibleOutputStruct() {
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

func (s *testYAMLLoaderSuite) TestLoad() {
	// Create a YAML file to read from
	filename := "test.yaml"
	if err := os.WriteFile(filename, []byte(testYAML), os.FileMode(0700)); err != nil {
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

func (s *testYAMLLoaderSuite) TestTrimExt() {
	s.Run("yaml gets trimmed", func() {
		filename := "test.yaml"
		name := s.loader.TrimExt(filename)

		s.Assert().Equal("test", name)
	})

	s.Run("yml gets trimmed", func() {
		filename := "test.yml"
		name := s.loader.TrimExt(filename)

		s.Assert().Equal("test", name)
	})

	s.Run("json does not get trimmed", func() {
		filename := "test.json"
		name := s.loader.TrimExt(filename)

		s.Assert().Equal("test.json", name)
	})
}
