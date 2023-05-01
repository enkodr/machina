package vm

import (
	"fmt"
	"os"
	"strings"

	"github.com/enkodr/machina/internal/netutil"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

var endpoint = "https://raw.githubusercontent.com/enkodr/machina/main/templates"

type Filer interface {
	Load() (*VMConfig, error)
}

type LocalFile struct {
	path string
}
type RemoteFile struct {
	name string
}

func NewTemplate(name string) Filer {
	if strings.Contains(name, ".yaml") {
		return &LocalFile{path: name}
	} else {
		return &RemoteFile{name: name}
	}
}

func (f *LocalFile) Load() (*VMConfig, error) {
	// Get the template content
	tpl, err := os.ReadFile(f.path)

	if err != nil {
		return nil, err
	}

	// Parse the YAML to struct
	vm, err := parseTemplate(tpl)
	if err != nil {
		return nil, err
	}
	return vm, nil
}

func (f *RemoteFile) Load() (*VMConfig, error) {
	// Get the template content
	tplFile := fmt.Sprintf("%s/%s.yaml", endpoint, f.name)
	tpl, err := netutil.Download(tplFile)

	// Parse the YAML to struct
	vm, err := parseTemplate(tpl)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

// parse the template from yaml to struct
func parseTemplate(tpl []byte) (*VMConfig, error) {
	vm, err := parseYaml(tpl)

	if err != nil {
		return nil, err
	}

	if vm.Extends != "" {
		tplFile := fmt.Sprintf("%s/%s.yaml", endpoint, vm.Extends)
		baseTpl, err := netutil.Download(tplFile)
		if err != nil {
			return nil, err
		}

		base, err := parseYaml(baseTpl)
		if err != nil {
			return nil, err
		}

		mergo.Merge(vm, base)
	}

	return vm, nil
}

// parse the template from yaml to struct
func parseYaml(tpl []byte) (*VMConfig, error) {
	vm := &VMConfig{}

	err := yaml.Unmarshal(tpl, vm)
	if err != nil {
		return nil, err
	}

	return vm, nil
}
