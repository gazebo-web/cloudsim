package loader

import (
	"github.com/gazebo-web/cloudsim/pkg/utils/reflect"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrInvalidOutValue is returned when a load function is passed an invalid out parameter.
	ErrInvalidOutValue = errors.New("invalid out value")
	// ErrLoadFailed is returned when a file failed to be loaded.
	ErrLoadFailed = errors.New("failed to load file")
)

// Loader reads and parses files into Go structs
type Loader interface {
	// Load loads a file from a specific path into the passed output value.
	//	 Loaders that do not use files can ignore the `path` parameter.
	Load(path string, out interface{}) error
	// TrimExt returns the filename without the file extension.
	//	 Loaders should remove the extension of the file they expect as config.
	//	 - A YAML Loader will remove the `.yaml` extension.
	//	 - A JSON Loader will remove the `.json` extension.
	TrimExt(filename string) string
	// Filter returns the list of files that match the extension this Loader can process.
	Filter(list []string) []string
}

// LoadFile loads a single file into a target struct.
func LoadFile(loader Loader, path string, out interface{}) error {
	return loader.Load(path, out)
}

// LoadDirFiles loads files contained in a directory into a target slice.
// Only a single directory level is loaded. Subdirectories are ignored. Failure to read a file will eventually be
// returned as an error, but will not stop the function from loading the rest of files in the target directory.
// `path` must be a directory path.
// `out` must be a map. The map must have string keys and pointer to interface type the files will be placed in as
// values (e.g. map[string]*Target). The keys of the map will contain filenames of loaded files, and the values
// will contain the loaded files.
//
// The return value will contain all errors found when attempting to load files.
// Returned error types can be checked using `errors.Is`. Errors of type ErrLoadFailed indicate that the file was
// not accessible. It is up to the caller to consider this a critical error. All other errors are critical.
func LoadDirFiles(loader Loader, path string, out interface{}) []error {
	// Get the directory's list of files
	files, err := os.ReadDir(path)
	if err != nil {
		return []error{err}
	}

	// Process files
	errs := make([]error, 0)
	for _, file := range files {
		// Only process files
		if file.IsDir() {
			continue
		}

		// Get a value instance to load the target file into
		value, err := reflect.NewCollectionValueInstance(out)
		if err != nil {
			errs = append(errs, errors.Wrap(ErrInvalidOutValue, err.Error()))
			continue
		}

		// Load the target file
		err = LoadFile(loader, filepath.Join(path, file.Name()), value)
		if err != nil {
			errs = append(errs, errors.Wrap(ErrLoadFailed, err.Error()))
			continue
		}

		// Append the loaded file to the output
		err = reflect.SetMapValue(out, file.Name(), value)
		if err != nil {
			errs = append(errs, errors.Wrap(ErrInvalidOutValue, err.Error()))
			continue
		}
	}

	// Return errors if found
	if len(errs) > 0 {
		return errs
	}

	return nil
}

// trimExts attempts to remove the given extensions from the filename.
//
//	Input: (file.yaml, .yaml)
//	Output: (file)
func trimExts(filename string, exts ...string) string {
	for _, ext := range exts {
		filename = strings.TrimSuffix(filename, ext)
	}
	return filename
}
