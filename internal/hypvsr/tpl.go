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

// Templater is an interface for loading different types of templates
type Templater interface {
	// Load method is responsible for loading the template
	// and returning an instance of KindManager and error if any occurs during the loading
	Load() (KindManager, error)
}

// LocalTemplate represents a local file-based template
type LocalTemplate struct {
	path string // path is the file system path to the local template
	name string // name is the name of the local template
}

// RemoteTemplate represents a template that needs to be fetched from a remote source
type RemoteTemplate struct {
	name string // name is the name of the remote template
}

// kind is a struct that represents the kind field in a yaml file
type kind struct {
	Kind string `yaml:"kind"` // Kind is the kind field in a yaml file
}

// NewTemplate is a factory function that returns an instance of Templater.
// It determines the type of Templater (LocalTemplate or RemoteTemplate)
// based on whether a file with the given name exists on the local file system.
func NewTemplate(name string) Templater {
	// Check if the passed argument name is a path to an existing file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		// If the file does not exist, assume it is a remote template
		return &RemoteTemplate{name: name}
	} else {
		// If the file does exist, assume it is a local template
		return &LocalTemplate{path: name}
	}
}

// Load is a method on the LocalTemplate struct that implements the Templater interface.
// It reads the content of the template from the local file system, parses it, and returns a corresponding KindManager.
func (f *LocalTemplate) Load() (KindManager, error) {
	// Reads the file named by filename and returns the contents
	tpl, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	// Call the template parser
	// If parsing the template content returns an error, it is propagated up
	vm, err := parseTemplate(tpl)
	if err != nil {
		return nil, err
	}

	// If there's no error, return the KindManager and a nil error
	return vm, nil
}

// Load is a method on the RemoteTemplate struct that implements the Templater interface.
// It reads the content of the template from the local file system, parses it, and returns a corresponding KindManager.
func (f *RemoteTemplate) Load() (KindManager, error) {
	// Set the URL from where to download the file
	tplFile := fmt.Sprintf("%s/%s.yaml", endpoint, f.name)

	// Dowload the template file from the remote endpoint
	tpl, err := netutil.Download(tplFile)
	if err != nil {
		return nil, err
	}

	// Call the template parser
	// If parsing the template content returns an error, it is propagated up
	vm, err := parseTemplate(tpl)
	if err != nil {
		return nil, err
	}

	// If there's no error, return the KindManager and a nil error
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
	data, err := os.ReadFile(filepath.Join(cfg.Directories.Instances, name, config.GetFilename(config.InstanceFilename)))

	// Loads the YAML to identify the km
	var k kind
	err = yaml.Unmarshal(data, &k)
	if err != nil {
		return nil, err
	}

	var km KindManager
	switch k.Kind {
	case "Machine":
		km = &Machine{}
		// Unmarshal the Machine
		err = yaml.Unmarshal(data, km)
		if err != nil {
			return nil, err
		}
		break
	case "Cluster":
		km = &Cluster{}
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
		vm.baseDir = cfg.Directories.Instances

		return vm, nil
	case "Cluster":
		c := &Cluster{}
		err := yaml.Unmarshal(tpl, c)
		if err != nil {
			return nil, err
		}
		c.baseDir = filepath.Join(cfg.Directories.Clusters, c.Name)
		expandedMachines := []Machine{}
		for _, machine := range c.Instances {
			machine.extend()

			if machine.Replicas == 0 {
				machine.Replicas = 1
			}
			for i := 0; i < machine.Replicas; i++ {
				copiedMachine := machine

				if machine.Replicas > 1 {
					copiedMachine.Name = fmt.Sprintf("%s-%d", copiedMachine.Name, i+1)
				}
				copiedMachine.baseDir = c.baseDir
				expandedMachines = append(expandedMachines, copiedMachine)
			}
		}
		c.Instances = expandedMachines

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
