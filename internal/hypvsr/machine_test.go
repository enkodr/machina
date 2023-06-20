package hypvsr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/osutil"
	"github.com/stretchr/testify/assert"
)

func TestMachine_CreateDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// Initialize a Machine machine for testing
	machine := Machine{
		baseDir: tmpDir,
		Name:    "test-machine",
	}

	// Test case: Machine directory does not exist
	err := machine.CreateDir()
	assert.NoError(t, err, "Error creating machine directory")

	// Test case: Machine directory already exists
	err = machine.CreateDir()
	assert.Error(t, err, "Machine directory should already exist")
}

func TestMachine_Prepare(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	// Initialize a Machine machine for testing
	machine := Machine{
		Name:    "test-machine",
		baseDir: tempDir,
		Credentials: Credentials{
			Username: "test-user",
			Password: "test-password",
			Groups:   []string{"group1", "group2"},
		},
	}

	// Mock configuration
	cfg = &config.Config{
		Directories: config.Directories{
			Images:    filepath.Join(tempDir, "images"),
			Instances: filepath.Join(tempDir, "instances"),
		},
	}

	// Create the machine directory for testing
	os.Mkdir(filepath.Join(tempDir, machine.Name), 0755)

	err := machine.Prepare()
	assert.NoError(t, err, "Error preparing machine")

	// Verify network configuration file
	networkPath := filepath.Join(tempDir, machine.Name, config.GetFilename(config.NetworkFilename))
	_, err = os.Stat(networkPath)
	assert.NoError(t, err, "Network configuration file not found")

	// Verify user data file
	userdataPath := filepath.Join(tempDir, machine.Name, config.GetFilename(config.UserdataFilename))
	_, err = os.Stat(userdataPath)
	assert.NoError(t, err, "User data file not found")

	// Verify private key file
	privateKeyPath := filepath.Join(tempDir, machine.Name, config.GetFilename(config.PrivateKeyFilename))
	_, err = os.Stat(privateKeyPath)
	assert.NoError(t, err, "Private key file not found")

	// Verify machine file
	machinePath := filepath.Join(tempDir, machine.Name, config.GetFilename(config.InstanceFilename))
	_, err = os.Stat(machinePath)
	assert.NoError(t, err, "Machine file not found")
}

func TestMachine_DownloadImage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	// Create a test image file
	fileContent := []byte("mock image content")

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a test image file
		imageData := []byte(fileContent)
		w.Write(imageData)
	}))

	// Calculate the SHA256 hash of the image data
	hash := sha256.Sum256([]byte(fileContent))
	checksum := "sha256:" + hex.EncodeToString(hash[:])

	// Mock configuration
	cfg = &config.Config{
		Directories: config.Directories{
			Images: tempDir,
		},
	}

	// Initialize a Machine machine for testing
	machine := Machine{
		Image: Image{
			URL:      fmt.Sprintf("%s/file.txt", mockServer.URL),
			Checksum: checksum,
		},
	}

	// Test case: Image is already downloaded
	err := machine.DownloadImage()
	assert.NoError(t, err, "Error downloading image")

	// Test case: Image needs to be downloaded
	machine.Image.URL = mockServer.URL + "/new-image.qcow2"
	err = machine.DownloadImage()
	assert.NoError(t, err, "Error downloading new image")

	// Test case: Invalid image URL
	machine.Image.URL = "invalid-url"
	err = machine.DownloadImage()
	assert.Error(t, err, "Invalid image URL")

	// Close the mock server
	mockServer.Close()
}

// MockRunner is a mock implementation of the Runner interface.
type MockRunner struct {
	Command string
	Args    []string
	Options []osutil.Option
	Output  string
	Error   error
	Called  bool
}

// RunCommand is the implementation of the Runner interface for the mock.
func (m *MockRunner) RunCommand(command string, args []string, options ...osutil.Option) (string, error) {
	m.Called = true
	m.Command = command
	m.Args = args
	m.Options = options
	return m.Output, m.Error
}

func TestCreateMachineDisk(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// Create a test image file
	fileContent := []byte("mock image content")
	imagePath := filepath.Join(tmpDir, "image.img")
	os.WriteFile(imagePath, fileContent, 0644)

	// Create an machine of your struct
	machine := &Machine{
		Image: Image{
			URL: imagePath,
		},
		Runner:  &MockRunner{},
		baseDir: filepath.Join(tmpDir, "test-machine"),
		Resources: Resources{
			Disk: "10G",
		},
	}

	cfg = &config.Config{
		Directories: config.Directories{
			Images: tmpDir,
		},
	}

	err := machine.createInstanceDisk()

	// Verify the first RunCommand call to create the disk
	mockRunner := machine.Runner.(*MockRunner)
	expectedCommand := "qemu-img"
	expectedArgs := []string{
		"create",
		"-F", "qcow2",
		"-b", imagePath,
		"-f", "qcow2", filepath.Join(tmpDir, "test-machine", "disk.img"),
		"10G",
	}
	assert.True(t, mockRunner.Called, "RunCommand should have been called")
	assert.Equal(t, expectedCommand, mockRunner.Command, "Unexpected command")
	assert.Equal(t, expectedArgs, mockRunner.Args, "Unexpected arguments")
	assert.NoError(t, err, "Unexpected error")
}

func TestCreateSeedDisk(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// Create an machine of your struct
	machine := &Machine{
		Runner: &MockRunner{},
		// Set other necessary fields for the test
	}

	err := machine.createSeedDisk()

	// Verify the RunCommand call to create the seed disk
	mockRunner := machine.Runner.(*MockRunner)
	expectedCommand := "cloud-localds"
	expectedArgs := []string{
		"--network-config=network.cfg",
		"seed.img",
		"userdata.yaml",
	}
	assert.True(t, mockRunner.Called, "RunCommand should have been called")
	assert.Equal(t, expectedCommand, mockRunner.Command, "Unexpected command")
	assert.Equal(t, expectedArgs, mockRunner.Args, "Unexpected arguments")
	assert.NoError(t, err, "Unexpected error")
}
