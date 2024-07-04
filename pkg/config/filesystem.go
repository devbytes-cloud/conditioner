package config

import "os"

// Filesystem defines the operations available for interacting with the filesystem.
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
