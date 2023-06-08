package hypvsr

// Cluster holds the configuration details for a cluster of machines
type Cluster struct {
	Kind     string    `yaml:"kind"`     // Kind of the resource, should be 'Cluster'
	Name     string    `yaml:"name"`     // Name of the cluster. Must be unique in the system
	Machines []Machine `yaml:"machines"` // List of machines in the cluster
}

func (c *Cluster) Create() error {
	return nil
}

func (c *Cluster) Start() error {
	return nil
}

func (c *Cluster) Stop() error {
	return nil
}

func (c *Cluster) ForceStop() error {
	return nil
}

func (c *Cluster) Status() (string, error) {
	return "", nil
}

func (c *Cluster) Delete() error {
	return nil
}

func (c *Cluster) CopyContent(origin string, dest string) error {
	return nil
}

func (c *Cluster) Prepare() error {
	return nil
}

func (c *Cluster) DownloadImage() error {
	return nil
}

func (c *Cluster) CreateDisks() error {
	return nil
}

func (c *Cluster) Wait() error {
	return nil
}

func (c *Cluster) Shell() error {
	return nil
}

func (c *Cluster) RunInitScripts() error {
	return nil
}

func (c *Cluster) GetVMs() []Machine {
	return nil
}
