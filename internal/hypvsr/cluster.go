package hypvsr

import (
	"errors"
	"os"

	"github.com/enkodr/machina/internal/osutil"
)

// Cluster holds the configuration details for a cluster of machines
type Cluster struct {
	baseDir  string
	Kind     string    `yaml:"kind"`     // Kind of the resource, should be 'Cluster'
	Name     string    `yaml:"name"`     // Name of the cluster. Must be unique in the system
	Machines []Machine `yaml:"machines"` // List of machines in the cluster
}

func (c *Cluster) CreateDir() error {
	// Check if VM already exists
	_, err := os.Stat(c.baseDir)
	if !os.IsNotExist(err) {
		return errors.New("cluster already exists")
	}
	osutil.MkDir(c.baseDir)

	for _, vm := range c.Machines {
		vm.baseDir = c.baseDir
		err := vm.CreateDir()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) Prepare() error {
	// Run the prepare for each of the machines
	for _, vm := range c.Machines {
		err := vm.Prepare()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) DownloadImage() error {
	for _, vm := range c.Machines {
		err := vm.DownloadImage()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) CreateDisks() error {
	for _, vm := range c.Machines {
		err := vm.CreateDisks()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) Create() error {
	for _, vm := range c.Machines {
		err := vm.Create()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) Wait() error {
	for _, vm := range c.Machines {
		err := vm.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) Start() error {
	for _, vm := range c.Machines {
		err := vm.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) Stop() error {
	for _, vm := range c.Machines {
		err := vm.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) ForceStop() error {
	for _, vm := range c.Machines {
		err := vm.ForceStop()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) Status() (string, error) {

	return "", nil
}

func (c *Cluster) Delete() error {
	for _, vm := range c.Machines {
		err := vm.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) CopyContent(origin string, dest string) error {

	return nil
}

func (c *Cluster) RunInitScripts() error {
	for _, vm := range c.Machines {
		err := vm.RunInitScripts()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) Shell() error {
	return nil
}

func (c *Cluster) GetVMs() []Machine {
	return c.Machines
}
