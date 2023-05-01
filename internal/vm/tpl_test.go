package vm

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalFile_Load(t *testing.T) {
	// Create a temporary file
	file, err := os.CreateTemp("", "*.yaml")
	assert.Nil(t, err)
	defer os.Remove(file.Name())

	// Write some YAML to the file
	yaml := []byte("name: My Template\n")
	_, err = file.Write(yaml)
	assert.Nil(t, err)

	// Load the file using LocalFile.Load
	lf := &LocalFile{path: file.Name()}
	vm, err := lf.Load()
	assert.Nil(t, err)

	// Check that the loaded MachinaVM matches the original YAML
	assert.Equal(t, "My Template", vm.Name)
}

func TestRemoteFile_Load(t *testing.T) {
	// Start an HTTP server that serves some YAML
	yaml := []byte("name: My Template\n")
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(yaml)
	}))
	defer testServer.Close()

	// Load the remote file using RemoteFile.Load
	rf := &RemoteFile{name: "template"}
	endpoint = testServer.URL

	vm, err := rf.Load()
	assert.Nil(t, err)

	// Check that the loaded MachinaVM matches the original YAML
	assert.Equal(t, "My Template", vm.Name)
}

func TestParseYaml(t *testing.T) {
	// Parse some YAML into a MachinaVM struct
	yaml := []byte("name: My Template\n")
	vm, err := parseYaml(yaml)
	assert.Nil(t, err)

	// Check that the parsed MachinaVM matches the original YAML
	assert.Equal(t, "My Template", vm.Name)
}

func TestParseTemplate(t *testing.T) {
	// Parse some YAML into a MachinaVM struct
	yaml := []byte("name: My Template\n")
	vm, err := parseTemplate(yaml)
	assert.Nil(t, err)

	// Check that the parsed MachinaVM matches the original YAML
	assert.Equal(t, "My Template", vm.Name)
}
