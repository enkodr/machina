package hypvsr

import (
	"strconv"

	"github.com/enkodr/machina/internal/config"
)

var cfg *config.Config

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

func getHypervisor() Hypervisor {
	if cfg == nil {
		cfg, _ = config.LoadConfig()
	}
	if cfg.Hypervisor == "qemu" {
		return &Qemu{}
	} else {
		return &Libvirt{}
	}
}
