package imgutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/enkodr/machina/internal/config"
	"github.com/stretchr/testify/assert"
)

// Helper function to check if a directory exists
func directoryExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func TestEnsureDirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Test case when directories are successfully created
	cfg := &config.Config{
		Directories: config.Directories{
			Images:    filepath.Join(tempDir, "/path/to/images"),
			Instances: filepath.Join(tempDir, "/path/to/instances"),
		},
	}

	EnsureDirectories(cfg)
	assert.True(t, directoryExists(cfg.Directories.Images))
	assert.True(t, directoryExists(cfg.Directories.Instances))

	// Test case when directory creation fails
	cfgInvalid := &config.Config{
		Directories: config.Directories{
			Images:    "/invalid/path", // Invalid path
			Instances: "/invalid/instances",
		},
	}
	EnsureDirectories(cfg)
	assert.False(t, directoryExists(cfgInvalid.Directories.Images))
	assert.False(t, directoryExists(cfgInvalid.Directories.Instances))
}

func TestGetFilenameFromURL(t *testing.T) {
	// Test case with a valid URL
	url := "https://example.com/files/document.pdf"
	expectedFilename := "document.pdf"
	filename, err := GetFilenameFromURL(url)
	assert.Nil(t, err)
	assert.Equal(t, expectedFilename, filename)

	// Test case with a URL that doesn't contain a filename
	url = "https://example.com/files/"
	expectedError := "no file in the URL"
	filename, err = GetFilenameFromURL(url)
	assert.NotNil(t, err)
	assert.EqualError(t, err, expectedError)

	// Test case with a URL that has multiple segments
	url = "https://example.com/files/folder/document.pdf"
	expectedFilename = "document.pdf"
	filename, err = GetFilenameFromURL(url)
	assert.Nil(t, err)
	assert.Equal(t, expectedFilename, filename)

	// Test case with an empty URL
	url = ""
	expectedError = "empty URL"
	filename, err = GetFilenameFromURL(url)
	assert.NotNil(t, err)
	assert.EqualError(t, err, expectedError)

	// Test case with a URL that only contains the root
	url = "https://example.com"
	expectedError = "no file in the URL"
	filename, err = GetFilenameFromURL(url)
	assert.NotNil(t, err)
	assert.EqualError(t, err, expectedError)

	// Test case with an invalid URL format
	url = "example.com"
	expectedError = "invalid URL format"
	filename, err = GetFilenameFromURL(url)
	assert.NotNil(t, err)
	assert.EqualError(t, err, expectedError)
}
