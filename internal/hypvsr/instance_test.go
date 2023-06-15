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
	"github.com/stretchr/testify/assert"
)

func TestMachine_CreateDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// Initialize a Machine instance for testing
	vm := Machine{
		baseDir: tmpDir,
		Name:    "test-machine",
	}

	// Test case: Machine directory does not exist
	err := vm.CreateDir()
	assert.NoError(t, err, "Error creating machine directory")

	// Test case: Machine directory already exists
	err = vm.CreateDir()
	assert.Error(t, err, "Machine directory should already exist")
}

func TestMachine_Prepare(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	// Initialize a Machine instance for testing
	vm := Machine{
		Name:    "test-machine",
		baseDir: tempDir,
		Credentials: Credentials{
			Username: "test-user",
			Password: "test-password",
			Groups:   []string{"group1", "group2"},
		},
	}

	// Create the machine directory for testing
	os.Mkdir(filepath.Join(tempDir, vm.Name), 0755)

	err := vm.Prepare()
	assert.NoError(t, err, "Error preparing machine")

	// Verify network configuration file
	networkPath := filepath.Join(tempDir, vm.Name, config.GetFilename(config.NetworkFilename))
	_, err = os.Stat(networkPath)
	assert.NoError(t, err, "Network configuration file not found")

	// Verify user data file
	userdataPath := filepath.Join(tempDir, vm.Name, config.GetFilename(config.UserdataFilename))
	_, err = os.Stat(userdataPath)
	assert.NoError(t, err, "User data file not found")

	// Verify private key file
	privateKeyPath := filepath.Join(tempDir, vm.Name, config.GetFilename(config.PrivateKeyFilename))
	_, err = os.Stat(privateKeyPath)
	assert.NoError(t, err, "Private key file not found")

	// Verify machine file
	machinePath := filepath.Join(tempDir, vm.Name, config.GetFilename(config.InstanceFilename))
	_, err = os.Stat(machinePath)
	assert.NoError(t, err, "Machine file not found")
}

func TestMachine_DownloadImage(t *testing.T) {
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
			Images: t.TempDir(),
		},
	}

	// Initialize a Machine instance for testing
	vm := Machine{
		Image: Image{
			URL:      fmt.Sprintf("%s/file.txt", mockServer.URL),
			Checksum: checksum,
		},
	}

	// Test case: Image is already downloaded
	err := vm.DownloadImage()
	assert.NoError(t, err, "Error downloading image")

	// Test case: Image needs to be downloaded
	vm.Image.URL = mockServer.URL + "/new-image.qcow2"
	err = vm.DownloadImage()
	assert.NoError(t, err, "Error downloading new image")

	// Test case: Invalid image URL
	vm.Image.URL = "invalid-url"
	err = vm.DownloadImage()
	assert.Error(t, err, "Invalid image URL")

	// Close the mock server
	mockServer.Close()
}
