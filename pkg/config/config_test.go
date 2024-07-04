package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/devbytes-cloud/conditioner/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	mockFS := new(mocks.MockFilesystem)
	path, err := getPath()
	assert.NoError(t, err)

	t.Run("File Exists", func(t *testing.T) {
		mockFS.On("Stat", path).Return(nil, nil).Once()

		exists, err := Exists(mockFS)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Random Error", func(t *testing.T) {
		expectedError := errors.New("random error")
		mockFS.On("Stat", path).Return(nil, expectedError).Once()

		exists, err := Exists(mockFS)
		assert.Error(t, err, expectedError)
		assert.False(t, exists)
	})

	t.Run("File does not exist", func(t *testing.T) {
		mockFS.On("Stat", path).Return(nil, fmt.Errorf("%w", os.ErrNotExist)).Once()

		exists, err := Exists(mockFS)
		assert.Nil(t, err)
		assert.False(t, exists)
	})
}

func TestRead(t *testing.T) {
	mockFS := new(mocks.MockFilesystem)

	path, err := getPath()
	assert.NoError(t, err)

	t.Run("Error: reading config", func(t *testing.T) {
		expectedError := errors.New("not found")
		mockFS.On("ReadFile", path).Return(nil, expectedError).Once()

		cfg, err := Read(mockFS)

		assert.Error(t, err, expectedError)
		assert.Nil(t, cfg)
	})

	t.Run("Error: json unmarshall", func(t *testing.T) {
		mockFS.On("ReadFile", path).Return(nil, nil).Once()

		cfg, err := Read(mockFS)

		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("Success: json unmarshall", func(t *testing.T) {
		expectedConfig := &Config{
			WhoAmI:    false,
			AllowList: []string{"unit-test"},
		}

		byteArray, err := json.Marshal(expectedConfig)
		assert.NoError(t, err)

		mockFS.On("ReadFile", path).Return(byteArray, nil).Once()

		cfg, err := Read(mockFS)

		assert.NoError(t, err)
		assert.Equal(t, cfg, expectedConfig)
	})
}
