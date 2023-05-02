package config

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

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
