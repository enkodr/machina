package hypvsr

import (
	"strconv"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/imgutil"
)

// KindManager defines an interface for managing virtual machine (VM) or cluster instances.
// It provides a standard set of operations that can be performed on these instances, regardless of the underlying type.
type KindManager interface {
	Create() error                    // Create is responsible for creating a new instance.
	Start() error                     // Start initiates the instance.
	Stop() error                      // Stop attempts to gracefully stop the instance.
	ForceStop() error                 // ForceStop forcefully stops the instance.
	Status() (string, error)          // Status returns the current status of the instance.
	Delete() error                    // Delete removes the instance.
	CopyContent(string, string) error // CopyContent copies content from a source to a destination path within the instance.
	Prepare() error                   // Prepare gets the instance ready for use, typically involves steps like setting up the network, file systems, etc.
	DownloadImage() error             // DownloadImage downloads the necessary image for creating the instance.
	CreateDisks() error               // CreateDisks creates the necessary disks for the instance.
	Wait() error                      // Wait waits until the instance is ready to use.
	Shell() error                     // Shell provides an interactive shell to the instance.
	RunInitScripts() error            // RunInitScripts runs initialization scripts on the instance.
	GetVMs() []Instance               // GetVMs returns a slice of Machine that represents the VMs under management.
	CreateDir() error                 // CreateDir creates a necessary directory for the instance.
}

var cfg *config.Config

// NewInstance is a factory function that returns an instance of KindManager.
func NewInstance(template Templater) (KindManager, error) {
	// Instantiate a new template
	// Load the template
	kind, err := template.Load()
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

	return kind, nil
}

// convertMemory is a function that converts the template memory to a value used by the hypervisor
func convertMemory(memory string) (string, error) {
	ram := memory

	// Check if the memory is in GB or MB
	switch suffix := memory[len(memory)-1]; suffix {
	// If the memory is in GB, convert it to MB
	case 'G':
		mem, err := strconv.Atoi(memory[0 : len(memory)-1])
		if err != nil {
			return "", err
		}
		bytes := mem * 1024
		ram = strconv.Itoa(bytes)
	// If the memory is in MB, leave it as it is
	case 'M':
		mem, err := strconv.Atoi(memory[0 : len(memory)-1])
		if err != nil {
			return "", err
		}
		ram = strconv.Itoa(mem)
	}

	return ram, nil
}
