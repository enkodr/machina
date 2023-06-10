package config

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Filename int

const (
	NetworkFilename Filename = iota
	UserdataFilename
	PrivateKeyFilename
	MachineFilename
	SeedImageFilename
	DiskFilename
	PIDFilename
)

func GetFilename(fn Filename) string {
	switch fn {
	case NetworkFilename:
		return "network.cfg"
	case UserdataFilename:
		return "userdata.yaml"
	case PrivateKeyFilename:
		return "id_rsa"
	case MachineFilename:
		return "machine.yaml"
	case SeedImageFilename:
		return "seed.img"
	case DiskFilename:
		return "disk.img"
	case PIDFilename:
		return "vm.pid"
	}

	return ""
}

type Config struct {
	Hypervisor  string      `yaml:"hypervisor,omitempty"`
	Instances   string      `yaml:"instances,omitempty"`
	Connection  string      `yaml:"connection,omitempty"`
	Directories Directories `yaml:"directories,omitempty"`
}

type Directories struct {
	Images   string `yaml:"images,omitempty"`
	Machines string `yaml:"instances,omitempty"`
	Clusters string `yaml:"clusters,omitempty"`
}

var (
	baseDir = ".local/share/machina"
	cfgDir  = ".config/machina"
)

// LoadConfig loads the configuration from file or creates a new if it's not yet created
func LoadConfig() (*Config, error) {
	// User home directory path
	home, _ := os.UserHomeDir()
	// Config file path
	cfgFile := filepath.Join(home, cfgDir, "config.yaml")

	// The configuration to be used
	cfg := &Config{
		Hypervisor: getHypervisor(),
		Connection: getConnection(),
		Directories: Directories{
			Images:   filepath.Join(home, baseDir, "images"),
			Machines: filepath.Join(home, baseDir, "instances/isolated"),
			Clusters: filepath.Join(home, baseDir, "instances/clusters"),
		},
	}

	cfgBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		yamlData, _ := yaml.Marshal(cfg)
		os.WriteFile(cfgFile, yamlData, 0644)
	} else {
		err = yaml.Unmarshal(cfgBytes, cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func getHypervisor() string {
	hypervisor := "libvirt"

	if runtime.GOOS == "darwin" {
		hypervisor = "qemu"
	}

	return hypervisor
}

func getConnection() string {
	conenction := "qemu:///system"

	if runtime.GOOS == "darwin" {
		conenction = ""
	}

	return conenction
}
