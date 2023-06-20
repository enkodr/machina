package hypvsr

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/enkodr/machina/internal/config"
)

// Cluster holds the configuration details for a cluster of machines
type Cluster struct {
	Kind      string            `yaml:"kind"`     // Kind of the resource, should be 'Cluster'
	Name      string            `yaml:"name"`     // Name of the cluster. Must be unique in the system
	Params    map[string]string `yaml:"params"`   // Parameters for the cluster
	Instances []Machine         `yaml:"machines"` // List of machines in the cluster
	Results   []string          `yaml:"results"`
}

func (c *Cluster) Prepare() error {
	err := c.parseParams()
	if err != nil {
		return err
	}

	err = c.createOutputDir()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cluster) createOutputDir() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(cfg.Directories.Results, c.Name), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cluster) parseParams() error {
	for i, vm := range c.Instances {
		// Create a new template and parse the template string
		tmpl, err := template.New("").Parse(vm.Scripts.Install)
		if err != nil {
			return err
		}

		// Create a buffer to store the parsed template
		var buf bytes.Buffer
		err = executeTemplateWithFallback(tmpl, &buf, c)
		if err != nil {
			return err
		}

		// Set the parsed template to the instance
		c.Instances[i].Scripts.Install = buf.String()
	}

	c.parseResults()

	return nil
}

func executeTemplateWithFallback(tmpl *template.Template, buf *bytes.Buffer, data interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			// Handle the panic and return an error
			buf.Reset()
			err := fmt.Errorf("template execution panicked: %v", r)
			fmt.Println(err)
		}
	}()

	err := tmpl.Execute(buf, data)
	return err
}

// Parses the results
func (cluster *Cluster) parseResults() {
	for i, machine := range cluster.Instances {
		for _, output := range cluster.Results {
			cluster.Instances[i].Scripts.Install = strings.ReplaceAll(machine.Scripts.Install, fmt.Sprintf("$(results.%s)", output), fmt.Sprintf("/etc/machina/results/%s/%s", cluster.Name, output))
		}
	}
}
