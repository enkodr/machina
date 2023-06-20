package hypvsr

import (
	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/imgutil"
)

type Instance struct {
	Kind       string
	Machines   []Machine
	Config     config.Config
	Hypervisor Hypervisor
}

// NewInstance is a factory function that returns an instance of KindManager.
func NewInstance(template Templater) (*Instance, error) {
	// Instantiate a new template
	// Load the template
	instance, err := template.Load()
	if err != nil {
		return nil, err
	}

	// Load configuration
	if cfg == nil {
		cfg, err = config.LoadConfig()
	}

	if err != nil {
		return nil, err
	}

	// Ensure directories
	imgutil.EnsureDirectories(cfg)

	return instance, nil
}
