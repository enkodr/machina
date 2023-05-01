package imgutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/enkodr/machina/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestEnsureDirectories(t *testing.T) {
	name := "test"
	EnsureDirectories(name)
	cfg := config.LoadConfig()

	// Check if directories were created
	_, err := os.Stat(filepath.Join(cfg.Directories.Images))
	assert.Nil(t, err)

	_, err = os.Stat(filepath.Join(cfg.Directories.Instances, name))
	assert.Nil(t, err)

	// Cleanup test directory
	os.Remove(filepath.Join(cfg.Directories.Instances, name))
}

func TestGetFilenameFromURL(t *testing.T) {
	// Test case 1: Valid URL
	want := "name.img"
	got, _ := GetFilenameFromURL("https://www.linux.com/name.img")
	assert.Equal(t, want, got)

	// Test case 2: Invalid URL
	_, err := GetFilenameFromURL("https://www.linux.com")
	assert.NotNil(t, err)
}
