package hypvsr

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplate(t *testing.T) {
	testCases := []struct {
		name     string
		wantType Templater
	}{
		{"local.yaml", &LocalTemplate{}},
		{"default", &RemoteTemplate{}},
		{"ubuntu.yaml", &LocalTemplate{}},
		{"ubuntu", &RemoteTemplate{}},
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

// func TestLocalFileLoadValidData(t *testing.T) {
// 	want := []byte(`name: TestVM
// specs:
//   cpus: 2
//   memory: "2G"
//   disk: "50G"
// `)

// 	name := "template.yaml"
// 	f, _ := os.CreateTemp("", name)
// 	defer f.Close()
// 	defer os.Remove(f.Name())
// 	f.Write(want)
// 	localFile := &LocalTemplate{path: f.Name()}
// 	got, _ := localFile.Load()
// 	assert.Equal(t, got.Name, "TestVM")
// 	assert.Equal(t, got.Resources.CPUs, "2")
// 	assert.Equal(t, got.Resources.Memory, "2G")
// 	assert.Equal(t, got.Resources.Disk, "50G")
// }

func TestLocalFileLoadInvalidData(t *testing.T) {
	want := "invalid data"
	name := "template.yaml"
	f, _ := os.CreateTemp("", name)
	f.Write([]byte(want))
	localFile := &LocalTemplate{path: f.Name()}
	vm, err := localFile.Load()
	assert.Error(t, err)
	assert.Nil(t, vm)
	defer f.Close()
	defer os.Remove(f.Name())
}

// func TestRemoteFileValidName(t *testing.T) {
// 	name := "ubuntu"
// 	data, _ := os.ReadFile(fmt.Sprintf("../../templates/%s.yaml", name))
// 	want := &Machine{}
// 	err := yaml.Unmarshal(data, want)
// 	assert.NoError(t, err)

// 	remoteFile := &RemoteTemplate{name: name}
// 	got, _ := remoteFile.Load()

// 	assert.Equal(t, want.Name, got.Name)
// 	assert.Equal(t, want.Resources.CPUs, got.Resources.CPUs)
// 	assert.Equal(t, want.Resources.CPUs, got.Resources.CPUs)
// 	assert.Equal(t, want.Resources.Memory, got.Resources.Memory)
// 	assert.Equal(t, want.Resources.Disk, got.Resources.Disk)
// }

func TestRemoteFileInvalidName(t *testing.T) {
	name := "invalid"

	remoteFile := &RemoteTemplate{name: name}
	vm, err := remoteFile.Load()

	assert.Error(t, err)
	assert.Nil(t, vm)
}

// func TestParseTemplateValidInput(t *testing.T) {
// 	want := []byte(`name: TestVM
// specs:
//   cpus: 2
//   memory: "2G"
//   disk: "50G"
// `)

// 	got, err := parseTemplate(want)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "TestVM", got.Name)
// 	assert.Equal(t, "2", got.Resources.CPUs)
// 	assert.Equal(t, "2G", got.Resources.Memory)
// }

func TestParseTemplate_InvalidYAML(t *testing.T) {
	want := []byte(`invalid data`)

	got, err := parseTemplate(want)
	assert.Error(t, err)
	assert.Nil(t, got)
}

// func TestParseTemplateExtends(t *testing.T) {
// 	want := []byte(`name: TestVM
// extends: ubuntu
// `)

// 	got, err := parseTemplate(want)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "TestVM", got.Name)
// 	assert.Equal(t, "2", got.Resources.CPUs)
// 	assert.Equal(t, "2G", got.Resources.Memory)
// }