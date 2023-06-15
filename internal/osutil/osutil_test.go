package osutil

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

// Helper function to calculate SHA512 checksum
func sha512Checksum(data []byte) string {
	hasher := sha512.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
