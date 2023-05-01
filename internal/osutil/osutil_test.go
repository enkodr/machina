package osutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMkDir(t *testing.T) {
	// Define test data
	testDir := "testdir"

	// Clean up test directory in case it already exists
	os.RemoveAll(testDir)

	// Create test directory
	MkDir(testDir)

	// Check if test directory exists
	_, err := os.Stat(testDir)
	assert.NoError(t, err)

	// Clean up test directory
	err = os.RemoveAll(testDir)
	assert.NoError(t, err)
}

func TestChecksum(t *testing.T) {
	// Test case 1: sha256 checksum
	path := "testdata/file.txt"
	checksum := "sha256:5bee30d1de847f564feaeb1f8ad30c2e9ace4766b3fa8a9fa11be2b2f0cea2f4"
	want := true
	got := Checksum(path, checksum)
	assert.Equal(t, want, got)

	// Test case 2: sha512 checksum
	checksum = "sha512:9c5f698667819d664c5fce6e6bed8cc0c96c488b148377a9a88e7dd670a401ddfd9e746e2c7e9935617815cf103143dada72b8b25a075d2747c105528cd9acf3"
	want = true
	got = Checksum(path, checksum)
	assert.Equal(t, want, got)

	// Test case 3: invalid checksum
	checksum = "md5:invalidchecksum"
	want = false
	got = Checksum(path, checksum)
	assert.Equal(t, want, got)

	// Test case 4: invalid file path
	path = "invalidpath"
	checksum = "sha256:5bee30d1de847f564feaeb1f8ad30c2e9ace4766b3fa8a9fa11be2b2f0cea2f4"
	want = false
	got = Checksum(path, checksum)
	assert.Equal(t, want, got)

	// Test case 5: test empty sha
	path = "testdata/file.txt"
	want = false
	got = Checksum(path, "")
	assert.Equal(t, want, got)
}
