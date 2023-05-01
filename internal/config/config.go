package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Hypervisor  string      `yaml:"hypervisor"`
	Connection  string      `yaml:"connection"`
	Directories Directories `yaml:"directories"`
}

type Directories struct {
	Images    string `yaml:"images"`
	Instances string `yaml:"instances"`
}

var (
	baseDir    = ".local/share/machina"
	cfgDir     = ".config/machina"
	hypervisor = "libvirt"
)

// LoadConfig loads the configuration from file or creates a new if it's not yet created
func LoadConfig() *Config {
	// User home directory path
	home, _ := os.UserHomeDir()
	// Config file path
	cfgFile := filepath.Join(home, cfgDir, "config.yaml")

	// The configuration to be used
	cfg := &Config{
		Hypervisor: hypervisor,
		Connection: "qemu:///system",
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
