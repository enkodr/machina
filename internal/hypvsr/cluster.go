package hypvsr

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/enkodr/machina/internal/config"
)

// Cluster holds the configuration details for a cluster of machines
type Cluster struct {
	config   *config.Config
	Kind     string    `yaml:"kind"`     // Kind of the resource, should be 'Cluster'
	Name     string    `yaml:"name"`     // Name of the cluster. Must be unique in the system
	Machines []Machine `yaml:"machines"` // List of machines in the cluster
}

func (c *Cluster) Prepare() error {
	return nil
}

func (c *Cluster) CreateDir() error {
	// Check if VM already exists
	_, err := os.Stat(filepath.Join(c.config.Directories.Clusters, c.Name))
	if !os.IsNotExist(err) {
		return errors.New("cluster already exists")
	}

	return nil
}

func (c *Cluster) Create() error {
	for _, vm := range c.Machines {
		vm.Create()
	}
	return nil
}

func (c *Cluster) Start() error {
	for _, vm := range c.Machines {
		vm.Start()
	}
	return nil
}

func (c *Cluster) Stop() error {
	for _, vm := range c.Machines {
		vm.Stop()
	}
	return nil
}

func (c *Cluster) ForceStop() error {
	for _, vm := range c.Machines {
		vm.ForceStop()
	}
	return nil
}

func (c *Cluster) Status() (string, error) {
	for _, vm := range c.Machines {
		vm.Status()
	}
	return "", nil
}

func (c *Cluster) Delete() error {
	for _, vm := range c.Machines {
		vm.Delete()
	}
	return nil
}

func (c *Cluster) CopyContent(origin string, dest string) error {
	for _, vm := range c.Machines {
		vm.CopyContent(origin, dest)
	}
	return nil
}

func (c *Cluster) DownloadImage() error {
	for _, vm := range c.Machines {
		vm.DownloadImage()
	}
	return nil
}

func (c *Cluster) CreateDisks() error {
	for _, vm := range c.Machines {
		vm.CreateDisks()
	}
	return nil
}

func (c *Cluster) Wait() error {
	for _, vm := range c.Machines {
		vm.Wait()
	}
	return nil
}

func (c *Cluster) Shell() error {
	for _, vm := range c.Machines {
		vm.Shell()
	}
	return nil
}

func (c *Cluster) RunInitScripts() error {
	for _, vm := range c.Machines {
		vm.RunInitScripts()
	}
	return nil
}

func (c *Cluster) GetVMs() []Machine {
	return c.Machines
}
