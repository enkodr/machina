package vm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/enkodr/machina/internal/netutil"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

var endpoint = "https://raw.githubusercontent.com/enkodr/machina/main/templates"

type Templater interface {
	Load() (*VMConfig, error)
}

type LocalTemplate struct {
	path string
}
type RemoteTemplate struct {
	name string
}

func NewTemplate(name string) Templater {
	if strings.Contains(name, ".yaml") {
		return &LocalTemplate{path: name}
	} else {
		return &RemoteTemplate{name: name}
	}
}

func (f *LocalTemplate) Load() (*VMConfig, error) {
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

func (f *RemoteTemplate) Load() (*VMConfig, error) {
	// Get the template content
	tplFile := fmt.Sprintf("%s/%s.yaml", endpoint, f.name)
	tpl, err := netutil.Download(tplFile)
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

// parse the template from yaml to struct
func parseTemplate(tpl []byte) (*VMConfig, error) {
	vm := &VMConfig{}

	err := yaml.Unmarshal(tpl, vm)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for vm.Extends != "" {
		tplFile := fmt.Sprintf("%s/%s.yaml", endpoint, vm.Extends)
		baseTpl, err := netutil.Download(tplFile)
		if err != nil {
			return nil, err
		}

		base := &VMConfig{}
		err = yaml.Unmarshal(baseTpl, base)
		if err != nil {
			return nil, err
		}
		vm.Extends = base.Extends
		base.Scripts = Scripts{}
		base.Mount = Mount{}
		mergo.Merge(vm, base)
	}
	vm.Resources.Disk = strings.ToUpper(vm.Resources.Disk)
	vm.Resources.Memory = strings.ToUpper(vm.Resources.Memory)

	return vm, nil
}

// GetTemplateList gets the list of available templates
func GetTemplateList() []string {
	type GitHubContent struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	url := fmt.Sprintf("https://api.github.com/repos/enkodr/machina/contents/templates")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var contents []GitHubContent
	if err := json.Unmarshal(body, &contents); err != nil {
		return nil
	}

	var files []string
	for _, c := range contents {
		if c.Type == "file" {
			file := strings.Split(c.Name, ".")[0]
			files = append(files, file)
		}
	}

	return files
}

func GetTemplate(name string) (string, error) {
	tpl, err := netutil.Download(fmt.Sprintf("%s/%s.yaml", endpoint, name))
	if err != nil {
		return "", err
	}

	return string(tpl), err
}
