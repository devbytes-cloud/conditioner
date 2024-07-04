package config

import "os"

// Filesystem defines the operations available for interacting with the filesystem.
// It provides an abstraction over the standard file operations, allowing for easier testing and mocking.
type Filesystem interface {
	// ReadFile reads the file named by `name` and returns the contents.
	// It returns the file contents as a byte slice and any error encountered.
	ReadFile(name string) ([]byte, error)

	// Stat returns the FileInfo structure describing the file.
	// If there is an error, it will be of type *PathError.
	Stat(name string) (os.FileInfo, error)

	// WriteFile writes data to a file named by `name`.
	// If the file does not exist, WriteFile creates it with permissions `perm`;
	// otherwise WriteFile truncates it before writing.
	WriteFile(name string, data []byte, perm os.FileMode) error
}

// FS implements the Filesystem interface using the os package.
type FS struct{}

// ReadFile reads the file named by `name` and returns the contents.
// It leverages os.ReadFile to read the file and return its contents along with any error encountered.
func (f FS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Stat returns the FileInfo structure describing the file.
// It uses os.Stat to obtain the FileInfo of the specified file, returning any errors encountered.
func (f FS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// WriteFile writes data to a file named by `name`.
// It creates or truncates the file named by `name`, writing the provided data with the specified permissions.
// It uses os.WriteFile to perform the operation, returning any errors encountered.
func (f FS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}
