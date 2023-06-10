package hypvsr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/netutil"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

var endpoint = "https://raw.githubusercontent.com/enkodr/machina/main/templates"

type Templater interface {
	Load() (KindManager, error)
}

type LocalTemplate struct {
	path string
	name string
}
type RemoteTemplate struct {
	name string
}

type kind struct {
	Kind string `yaml:"kind"`
}

func NewTemplate(name string) Templater {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return &RemoteTemplate{name: name}
	} else {
		return &LocalTemplate{path: name}
	}
}

func (f *LocalTemplate) Load() (KindManager, error) {
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

func (f *RemoteTemplate) Load() (KindManager, error) {
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

// Load loads the YAML file content into the struct
func Load(name string) (KindManager, error) {
	// Loads the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Reads the YAML file
	data, err := os.ReadFile(filepath.Join(cfg.Directories.Machines, name, config.GetFilename(config.MachineFilename)))

	// Loads the YAML to identify the km
	var k kind
	err = yaml.Unmarshal(data, &k)
	if err != nil {
		return nil, err
	}

	var km KindManager
	switch k.Kind {
	case "Machine":
		km = &Machine{
			config: cfg,
		}
		// Unmarshal the Machine
		err = yaml.Unmarshal(data, km)
		if err != nil {
			return nil, err
		}
		break
	case "Cluster":
		km = &Cluster{
			config: cfg,
		}
		// Unmarshal the Machine
		err = yaml.Unmarshal(data, km)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown kind")
	}

	return km, nil
}

// parse the template from yaml to struct
func parseTemplate(tpl []byte) (KindManager, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Identify the kind
	k := &kind{}
	yaml.Unmarshal(tpl, k)

	switch k.Kind {
	case "Machine":
		vm := &Machine{}
		err := yaml.Unmarshal(tpl, vm)
		if err != nil {
			return nil, err
		}
		err = vm.extend()
		if err != nil {
			return nil, err
		}
		vm.config = cfg
		return vm, nil
	case "Cluster":
		c := &Cluster{}
		err := yaml.Unmarshal(tpl, c)
		if err != nil {
			return nil, err
		}
		for _, vm := range c.Machines {
			vm.extend()
		}
		c.config = cfg
		return c, nil
	}

	return nil, errors.New("unsupported kind")
}

func (vm *Machine) extend() error {
	for vm.Extends != "" {
		tplFile := fmt.Sprintf("%s/%s.yaml", endpoint, vm.Extends)
		baseTpl, err := netutil.Download(tplFile)
		if err != nil {
			return err
		}

		base := &Machine{}
		err = yaml.Unmarshal(baseTpl, base)
		if err != nil {
			return err
		}
		vm.Extends = base.Extends
		base.Scripts = Scripts{}
		base.Mount = Mount{}
		mergo.Merge(vm, base)
	}
	vm.Resources.Disk = strings.ToUpper(vm.Resources.Disk)
	vm.Resources.Memory = strings.ToUpper(vm.Resources.Memory)

	return nil
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
