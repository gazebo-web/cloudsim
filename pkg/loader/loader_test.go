package loader

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type T struct {
	Value string
}

type file struct {
	Path  string
	Value string
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(testLoaderSuite))
}

type testLoaderSuite struct {
	suite.Suite
	tmpDir          string
	logger          ign.Logger
	loader          Loader
	applicationsDir string
	file1           file
	file2           file
	file3           file
}

func (s *testLoaderSuite) prefixPath(path string) string {
	return filepath.Join(s.tmpDir, path)
}

func (s *testLoaderSuite) createFile(path string, data string) error {
	path = s.prefixPath(path)
	return ioutil.WriteFile(path, []byte(data), 0777)
}

func (s *testLoaderSuite) createDir(path string) error {
	path = s.prefixPath(path)
	return os.Mkdir(path, 0700)
}

func (s *testLoaderSuite) deleteTmpDir() error {
	return os.RemoveAll(s.tmpDir)
}

func (s *testLoaderSuite) SetupSuite() {
	s.logger = ign.NewLoggerNoRollbar("testLoaderSuite", ign.VerbosityWarning)
	// Using a YAML loader, but any loader will do
	s.loader = NewYAMLLoader(s.logger)

	// Create a temporary directory for test files
	var err error
	s.tmpDir, err = ioutil.TempDir("", "testLoaderSuite*")
	s.Equal(nil, err)

	// Create directories
	s.applicationsDir = "applications"
	s.Equal(nil, s.createDir(s.applicationsDir))
	s.Equal(nil, s.createDir(filepath.Join(s.applicationsDir, "subdir1")))
	s.Equal(nil, s.createDir(filepath.Join(s.applicationsDir, "subdir2")))

	// Create test files
	s.file1 = file{
		Path:  filepath.Join(s.applicationsDir, "file1"),
		Value: "1",
	}
	s.Equal(nil, s.createFile(s.file1.Path, "value: 1"))

	s.file2 = file{
		Path:  filepath.Join(s.applicationsDir, "file2"),
		Value: "2",
	}
	s.Equal(nil, s.createFile(s.file2.Path, "value: 2"))

	s.file3 = file{
		Path:  filepath.Join(s.applicationsDir, "file3"),
		Value: "",
	}
	s.Equal(nil, s.createFile(s.file3.Path, `a:-`))
}

func (s *testLoaderSuite) TearDownSuite() {
	// Delete the temporary directory and all files within it
	s.Equal(nil, s.deleteTmpDir())
}

func (s *testLoaderSuite) TestLoadFileWithFile() {
	// Load file
	path := s.prefixPath(s.file1.Path)
	out := T{}
	s.Equal(nil, LoadFile(s.loader, path, &out))

	// Verify file contents
	s.Equal(s.file1.Value, out.Value)
}

func (s *testLoaderSuite) TestLoadFileWithInvalidFile() {
	path := s.prefixPath(s.file3.Path)
	out := T{}
	s.NotEqual(nil, LoadFile(s.loader, path, &out))
}

func (s *testLoaderSuite) TestLoadFileWithDir() {
	// Should fail to load dir as file
	path := s.prefixPath(s.applicationsDir)
	out := T{}
	s.True(errors.Is(LoadFile(s.loader, path, &out), ErrLoadFailed))
}

func (s *testLoaderSuite) TestLoadDirFilesWithFile() {
	// Load file
	path := s.prefixPath(s.file1.Path)
	out := T{}
	s.Equal(nil, LoadFile(s.loader, path, &out))

	// Verify file contents
	s.Equal(s.file1.Value, out.Value)
}

func (s *testLoaderSuite) TestLoadDirFilesWithDir() {
	path := s.prefixPath(s.applicationsDir)

	// Load files
	out := make(map[string]*T, 0)
	errs := LoadDirFiles(s.loader, path, out)

	// The map should contain files 1 and 2 and fail to load file 3 because its syntax is invalid.
	s.Len(out, 2)
	_, ok := out[filepath.Base(s.file1.Path)]
	s.True(ok)
	_, ok = out[filepath.Base(s.file2.Path)]
	s.True(ok)

	// File 3 should have failed to load and be contained in the returned errors
	s.Len(errs, 1)
	s.True(errors.Is(errs[0], ErrLoadFailed))
}
