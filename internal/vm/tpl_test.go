package vm

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNewTemplate(t *testing.T) {
	testCases := []struct {
		name     string
		wantType Filer
	}{
		{"local.yaml", &LocalFile{}},
		{"default", &RemoteFile{}},
		{"ubuntu.yaml", &LocalFile{}},
		{"ubuntu", &RemoteFile{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewTemplate(tc.name)

			switch r := result.(type) {
			case any:
				assert.IsType(t, tc.wantType, r)
			}
		})
	}
}

func TestLocalFileLoadValidData(t *testing.T) {
	want := []byte(`name: TestVM
specs:
  cpus: 2
  memory: "2G"
  disk: "50G"
`)

	name := "template.yaml"
	f, _ := os.CreateTemp("", name)
	defer f.Close()
	defer os.Remove(f.Name())
	f.Write(want)
	localFile := &LocalFile{path: f.Name()}
	got, _ := localFile.Load()
	assert.Equal(t, got.Name, "TestVM")
	assert.Equal(t, got.Specs.CPUs, "2")
	assert.Equal(t, got.Specs.Memory, "2G")
	assert.Equal(t, got.Specs.Disk, "50G")
}

func TestLocalFileLoadInvalidData(t *testing.T) {
	want := "invalid data"
	name := "template.yaml"
	f, _ := os.CreateTemp("", name)
	f.Write([]byte(want))
	localFile := &LocalFile{path: f.Name()}
	vm, err := localFile.Load()
	assert.Error(t, err)
	assert.Nil(t, vm)
	defer f.Close()
	defer os.Remove(f.Name())
}

func TestRemoteFileValidName(t *testing.T) {
	name := "ubuntu"
	data, _ := os.ReadFile(fmt.Sprintf("../../templates/%s.yaml", name))
	want := &VMConfig{}
	err := yaml.Unmarshal(data, want)
	assert.NoError(t, err)

	remoteFile := &RemoteFile{name: name}
	got, _ := remoteFile.Load()

	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Specs.CPUs, got.Specs.CPUs)
	assert.Equal(t, want.Image.Checksum, got.Image.Checksum)

}

func TestRemoteFileInvalidName(t *testing.T) {
	name := "invalid"

	remoteFile := &RemoteFile{name: name}
	vm, err := remoteFile.Load()

	assert.Error(t, err)
	assert.Nil(t, vm)
}

func TestParseTemplateValidInput(t *testing.T) {
	want := []byte(`name: TestVM
specs:
  cpus: 2
  memory: "2G"
  disk: "50G"
`)

	got, err := parseTemplate(want)
	assert.NoError(t, err)
	assert.Equal(t, "TestVM", got.Name)
	assert.Equal(t, "2", got.Specs.CPUs)
	assert.Equal(t, "2G", got.Specs.Memory)
}

func TestParseTemplate_InvalidYAML(t *testing.T) {
	want := []byte(`invalid data`)

	got, err := parseTemplate(want)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestParseTemplateExtends(t *testing.T) {
	want := []byte(`name: TestVM
extends: default
`)

	got, err := parseTemplate(want)
	assert.NoError(t, err)
	assert.Equal(t, "TestVM", got.Name)
	assert.Equal(t, "2", got.Specs.CPUs)
	assert.Equal(t, "2G", got.Specs.Memory)
}
