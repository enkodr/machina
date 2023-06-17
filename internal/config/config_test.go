package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGetConfigFilePath(t *testing.T) {
	// Set up a temporary home directory for testing
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	expectedPath := filepath.Join(tempDir, cfgDir, "config.yaml")
	actualPath := getConfigFilePath()

	assert.Equal(t, expectedPath, actualPath, "Incorrect config file path")
}

func TestConfigExists_FileExists(t *testing.T) {
	// Create a temporary config file for testing
	tempFile, err := ioutil.TempFile("", "config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	exists := configExists(tempFile.Name())

	assert.True(t, exists, "Expected config file to exist")
}

func TestConfigExists_FileNotExists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.yaml")

	exists := configExists(filePath)

	assert.False(t, exists, "Expected config file to not exist")
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file with sample content
	tempFile, err := ioutil.TempFile("", "config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	sampleConfig := `
hypervisor: test-hypervisor
connection: test-connection
directories:
  images: /path/to/images
  instances: /path/to/instances
  clusters: /path/to/clusters
`
	err = ioutil.WriteFile(tempFile.Name(), []byte(sampleConfig), 0644)
	assert.NoError(t, err)

	expectedConfig := &Config{
		Hypervisor:  "test-hypervisor",
		Connection:  "test-connection",
		Directories: Directories{Images: "/path/to/images", Instances: "/path/to/instances"},
	}

	actualConfig, err := loadConfigFromFile(tempFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, expectedConfig, actualConfig, "Loaded config does not match expected config")
}

func TestCreateDefaultConfig(t *testing.T) {
	// Set up a temporary home directory for testing
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Set up a temporary config directory for testing
	configDir := filepath.Join(tempDir, cfgDir)
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	tempFile := filepath.Join(configDir, "config.yaml")

	expectedConfig := &Config{
		Hypervisor: getHypervisor(),
		Connection: getConnection(),
		Directories: Directories{
			Images:    getDefaultImagePath(),
			Instances: getDefaultInstancesPath(),
		},
	}

	_, err = createDefaultConfig(tempFile)
	assert.NoError(t, err)

	actualBytes, err := ioutil.ReadFile(tempFile)
	assert.NoError(t, err)

	actualConfig := &Config{}
	err = yaml.Unmarshal(actualBytes, actualConfig)
	assert.NoError(t, err)

	assert.Equal(t, expectedConfig, actualConfig, "Created config does not match expected config")
}
