package hypvsr

import (
	"errors"
	"os"
)

// Cluster holds the configuration details for a cluster of machines
type Cluster struct {
	baseDir   string
	Kind      string     `yaml:"kind"`     // Kind of the resource, should be 'Cluster'
	Name      string     `yaml:"name"`     // Name of the cluster. Must be unique in the system
	Instances []Instance `yaml:"machines"` // List of machines in the cluster
}

// CreateDir method creates the directory where the instance files will be stored
func (c *Cluster) CreateDir() error {
	// Check if VM already exists
	_, err := os.Stat(c.baseDir)
	if !os.IsNotExist(err) {
		return errors.New("cluster already exists")
	}

	// Create the cluster direcotry
	os.Mkdir(c.baseDir, 0755)

	// Create the directory for each of the cluster instances
	for _, vm := range c.Instances {
		vm.baseDir = c.baseDir
		err := vm.CreateDir()
		if err != nil {
			return err
		}
	}
	return nil
}

// Prepare method calls the Prepare method for each instance
func (c *Cluster) Prepare() error {
	// Run the prepare for each of the machines
	for _, vm := range c.Instances {
		err := vm.Prepare()
		if err != nil {
			return err
		}
	}
	return nil
}

// DownloadImage method calls the DownloadImage method for each instance
func (c *Cluster) DownloadImage() error {
	for _, vm := range c.Instances {
		err := vm.DownloadImage()
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDisks method calls the CreateDisks method for each instance
func (c *Cluster) CreateDisks() error {
	for _, vm := range c.Instances {
		err := vm.CreateDisks()
		if err != nil {
			return err
		}
	}
	return nil
}

// Create method calls the Create method for each instance
func (c *Cluster) Create() error {
	for _, vm := range c.Instances {
		err := vm.Create()
		if err != nil {
			return err
		}
	}
	return nil
}

// Wait method calls the Wait method for each instance
func (c *Cluster) Wait() error {
	for _, vm := range c.Instances {
		err := vm.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

// Start method calls the Start method for each instance
func (c *Cluster) Start() error {
	for _, vm := range c.Instances {
		err := vm.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop method calls the Stop method for each instance
func (c *Cluster) Stop() error {
	for _, vm := range c.Instances {
		err := vm.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

// ForceStop method calls the ForceStop method for each instance
func (c *Cluster) ForceStop() error {
	for _, vm := range c.Instances {
		err := vm.ForceStop()
		if err != nil {
			return err
		}
	}
	return nil
}

// Status method calls the Status method for each instance
func (c *Cluster) Status() (string, error) {

	return "", nil
}

// Delete method calls the Delete method for each instance
func (c *Cluster) Delete() error {
	for _, vm := range c.Instances {
		err := vm.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyContent method calls the CopyContent method for each instance
func (c *Cluster) CopyContent(origin string, dest string) error {

	return nil
}

// RunInitScripts method calls the RunInitScripts method for each instance
func (c *Cluster) RunInitScripts() error {
	for _, vm := range c.Instances {
		err := vm.RunInitScripts()
		if err != nil {
			return err
		}
	}
	return nil
}

// Shell method calls the Shell method for each instance
func (c *Cluster) Shell() error {
	return nil
}

// GetVMs method calls the GetVMs method for each instance
func (c *Cluster) GetVMs() []Instance {
	return c.Instances
}
