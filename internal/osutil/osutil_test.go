package osutil

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandRunner_RunCommand(t *testing.T) {
	runner := CommandRunner{}
	command := "echo"
	args := []string{"Hello", "World"}

	output, err := runner.RunCommand(command, args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedOutput := "Hello World\n"
	if output != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, output)
	}
}

func TestCommandRunner_RunNonExisingCommand(t *testing.T) {
	runner := CommandRunner{}
	command := "nonexistent"
	args := []string{}

	_, err := runner.RunCommand(command, args)

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestCommandRunner_RunCommandNonExistingFile(t *testing.T) {
	// Define the command to run
	runner := CommandRunner{}
	command := "ls"
	args := []string{"nonexistent-file"}

	// Run the command
	_, err := runner.RunCommand(command, args)

	// Assert that an error is returned
	assert.Error(t, err, "Expected an error, got nil")
}

func TestChecksum(t *testing.T) {
	// Create a temporary file for testing
	fileContent := []byte("Test file content")
	tmpFile, err := ioutil.TempFile("", "testfile")
	assert.Nil(t, err)
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath)
	defer tmpFile.Close()
	_, err = tmpFile.Write(fileContent)
	assert.Nil(t, err)

	// Test case with a matching SHA256 checksum
	checksum := "sha256:" + sha256Checksum(fileContent)
	result := Checksum(tmpFilePath, checksum)
	assert.True(t, result)

	// Test case with a non-matching SHA256 checksum
	checksum = "sha256:InvalidChecksum"
	result = Checksum(tmpFilePath, checksum)
	assert.False(t, result)

	// Test case with an unsupported algorithm (MD5)
	checksum = "md5:InvalidChecksum"
	result = Checksum(tmpFilePath, checksum)
	assert.False(t, result)

	// Test case with an invalid checksum format
	checksum = "InvalidChecksum"
	result = Checksum(tmpFilePath, checksum)
	assert.False(t, result)

	// Test case with a non-existing file
	nonExistingFilePath := "/path/to/nonexistingfile"
	checksum = "sha256:" + sha256Checksum(fileContent)
	result = Checksum(nonExistingFilePath, checksum)
	assert.False(t, result)

	// Test case with a matching SHA512 checksum
	checksum = "sha512:" + sha512Checksum(fileContent)
	result = Checksum(tmpFilePath, checksum)
	assert.True(t, result)

	// Test case with a non-matching SHA512 checksum
	checksum = "sha512:InvalidChecksum"
	result = Checksum(tmpFilePath, checksum)
	assert.False(t, result)
}

// Helper function to calculate SHA256 checksum
func sha256Checksum(data []byte) string {
	// Calculate the SHA256 hash of the image data
	hasher := sha256.New()
	hasher.Write(data)

	// Return the SHA256 checksum
	return hex.EncodeToString(hasher.Sum(nil))
}

// Helper function to calculate SHA512 checksum
func sha512Checksum(data []byte) string {
	// Calculate the SHA512 hash of the image data
	hasher := sha512.New()
	hasher.Write(data)

	// Return the SHA512 checksum
	return hex.EncodeToString(hasher.Sum(nil))
}

// MockRunner is a mock implementation of the Runner interface.
type MockRunner struct {
	Command string
	Args    []string
	Options []Option
	Output  string
	Error   error
	Called  bool
}

// RunCommand is the implementation of the Runner interface for the mock.
func (m *MockRunner) RunCommand(command string, args []string, options ...Option) (string, error) {
	m.Called = true
	m.Command = command
	m.Args = args
	m.Options = options
	return m.Output, m.Error
}

func TestChecksumWithValidChecksum(t *testing.T) {
	// Create a temporary file for testing
	tmpFile := createTempFile(t, "Hello, World!")

	// Calculate the SHA256 checksum of the file
	hash := sha256.New()
	hash.Write([]byte("Hello, World!"))
	checksum := "sha256:" + hex.EncodeToString(hash.Sum(nil))

	// Verify the checksum
	result := Checksum(tmpFile, checksum)

	// Assert that the result is true
	assert.True(t, result, "Expected checksum verification to pass")
}

func TestChecksumWithInvalidChecksum(t *testing.T) {
	// Create a temporary file for testing
	tmpFile := createTempFile(t, "Hello, World!")

	// Specify an incorrect checksum
	incorrectChecksum := "sha256:123456789"

	// Verify the checksum
	result := Checksum(tmpFile, incorrectChecksum)

	// Assert that the result is false
	assert.False(t, result, "Expected checksum verification to fail")
}

func createTempFile(t *testing.T, content string) string {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		t.Fatal("Failed to create temporary file:", err)
	}

	// defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write the content to the temporary file
	_, err = io.WriteString(tmpFile, content)
	assert.Nil(t, err)

	// Return the path to the temporary file
	return tmpFile.Name()
}
