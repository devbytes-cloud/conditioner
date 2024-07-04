package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/mitchellh/go-homedir"
)

// ConfigName is the name of the configuration file.
const configName string = "~/.conditioner.json"

// Config represents the conditioners configuration.
// It includes fields for user preferences and settings.
type Config struct {
	// WhoAmI indicates whether to prepend the user's identity to the output.
	WhoAmI bool `json:"prepend-whoami"`
	// AllowList is a list of allowed entities for the application.
	AllowList []string `json:"allow-list"`
}

// Exists checks if the configuration file exists.
// It returns a boolean indicating the existence of the file and any error encountered.
func Exists(fs Filesystem) (bool, error) {
	path, err := getPath()
	if err != nil {
		return false, err
	}

	if _, err := fs.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

// Write writes the default configuration to the configuration file.
// It returns any error encountered during the operation.
func Write() error {
	conf := &Config{
		WhoAmI:    false,
		AllowList: []string{},
	}

	confJson, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}

	fileName, err := getPath()
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, confJson, 0o644)
}

// Read is a placeholder function for reading the configuration file.
func Read(fs Filesystem) (*Config, error) {
	path, err := getPath()
	if err != nil {
		return nil, err
	}

	byteConfig, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(byteConfig, config); err != nil {
		return nil, err
	}
	return config, nil
}

// getPath expands the configuration file name to its full path.
// It returns the full path and any error encountered.
func getPath() (string, error) {
	return homedir.Expand(configName)
}
