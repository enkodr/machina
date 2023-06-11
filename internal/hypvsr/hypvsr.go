package hypvsr

import (
	"strconv"

	"github.com/enkodr/machina/internal/imgutil"
)

// KindManager defines an interface for managing virtual machine (VM) or cluster instances.
// It provides a standard set of operations that can be performed on these instances, regardless of the underlying type.
type KindManager interface {
	// Create is responsible for creating a new instance.
	Create() error
	// Start initiates the instance.
	Start() error
	// Stop attempts to gracefully stop the instance.
	Stop() error
	// ForceStop forcefully stops the instance.
	ForceStop() error
	// Status returns the current status of the instance.
	Status() (string, error)
	// Delete removes the instance.
	Delete() error
	// CopyContent copies content from a source to a destination path within the instance.
	CopyContent(string, string) error
	// Prepare gets the instance ready for use, typically involves steps like setting up the network, file systems, etc.
	Prepare() error
	// DownloadImage downloads the necessary image for creating the instance.
	DownloadImage() error
	// CreateDisks creates the necessary disks for the instance.
	CreateDisks() error
	// Wait waits until the instance is ready to use.
	Wait() error
	// Shell provides an interactive shell to the instance.
	Shell() error
	// RunInitScripts runs initialization scripts on the instance.
	RunInitScripts() error
	// GetVMs returns a slice of Machine that represents the VMs under management.
	GetVMs() []Machine
	// CreateDir creates a necessary directory for the instance.
	CreateDir() error
}

// NewInstance creates a new machine
func NewInstance(name string) (KindManager, error) {
	// Instantiate a new template
	// Load the template
	kind, err := NewTemplate(name).Load()
	if err != nil {
		return nil, err
	}

	// Create directory structure
	err = imgutil.EnsureDirectories()
	if err != nil {
		return nil, err
	}

	return kind, nil
}

func convertMemory(memory string) (string, error) {
	ram := memory
	if memory[len(memory)-1] == 'G' {
		mem, err := strconv.Atoi(memory[0 : len(memory)-1])
		if err != nil {
			return "", err
		}
		bytes := mem * 1024
		ram = strconv.Itoa(bytes)
	}
	return ram, nil
}
