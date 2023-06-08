package config

import (
	"fmt"
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

func GetFilename(name string, fn Filename) string {
	switch fn {
	case NetworkFilename:
		return fmt.Sprintf("%s-network.cfg", name)
	case UserdataFilename:
		return fmt.Sprintf("%s-userdata.yaml", name)
	case PrivateKeyFilename:
		return fmt.Sprintf("%s-id_rsa", name)
	case MachineFilename:
		return fmt.Sprintf("%s-machine.yaml", name)
	case SeedImageFilename:
		return fmt.Sprintf("%s-seed.img", name)
	case DiskFilename:
		return fmt.Sprintf("%s-disk.img", name)
	case PIDFilename:
		return fmt.Sprintf("%s-vm.pid", name)
	}

	return ""
}

type Config struct {
	Hypervisor  string      `yaml:"hypervisor,omitempty"`
	Connection  string      `yaml:"connection,omitempty"`
	Directories Directories `yaml:"directories,omitempty"`
}

type Directories struct {
	Images    string `yaml:"images,omitempty"`
	Instances string `yaml:"instances,omitempty"`
}

var (
	baseDir = ".local/share/machina"
	cfgDir  = ".config/machina"
)

// LoadConfig loads the configuration from file or creates a new if it's not yet created
func LoadConfig() *Config {
	// User home directory path
	home, _ := os.UserHomeDir()
	// Config file path
	cfgFile := filepath.Join(home, cfgDir, "config.yaml")

	// The configuration to be used
	cfg := &Config{
		Hypervisor: getHypervisor(),
		Connection: getConnection(),
		Directories: Directories{
			Images:    filepath.Join(home, baseDir, "images"),
			Instances: filepath.Join(home, baseDir, "instances"),
		},
	}

	cfgBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		yamlData, _ := yaml.Marshal(cfg)
		os.WriteFile(cfgFile, yamlData, 0644)
	} else {
		yaml.Unmarshal(cfgBytes, cfg)
	}

	return cfg
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
