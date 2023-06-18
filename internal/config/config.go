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
	InstanceFilename
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
	case InstanceFilename:
		return "instance.yaml"
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
	Connection  string      `yaml:"connection,omitempty"`
	Directories Directories `yaml:"directories,omitempty"`
}

type Directories struct {
	Images    string `yaml:"images,omitempty"`
	Instances string `yaml:"instances,omitempty"`
	Results   string `yaml:"results,omitempty"`
}

var (
	baseDir = ".local/share/machina"
	cfgDir  = ".config/machina"
)

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	cfgFile := getConfigFilePath()

	// If the config file exists, load it
	if configExists(cfgFile) {
		return loadConfigFromFile(cfgFile)
	}

	// Otherwise, create a default config file
	return createDefaultConfig(cfgFile)
}

// GetConfigFilePath returns the path to the config file
func getConfigFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, cfgDir, "config.yaml")
}

// configExists checks if the config file exists
func configExists(cfgFile string) bool {
	_, err := os.Stat(cfgFile)
	return err == nil
}

// loadConfigFromFile loads the configuration from the config file
func loadConfigFromFile(cfgFile string) (*Config, error) {
	// Read the config file
	cfgBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal the config file
	var cfg Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// createDefaultConfig creates a default config file
func createDefaultConfig(cfgFile string) (*Config, error) {
	cfg := Config{
		Hypervisor: getHypervisor(),
		Connection: getConnection(),
		Directories: Directories{
			Images:    getDefaultImagePath(),
			Instances: getDefaultInstancesPath(),
			Results:   getDefaultResultsPath(),
		},
	}

	// Create the config file
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	// Write the config file
	err = os.WriteFile(cfgFile, yamlData, 0644)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetDefaultImagePath returns the default path for images
func getDefaultImagePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, baseDir, "images")
}

// GetDefaultInstancesPath returns the default path for instances
func getDefaultInstancesPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, baseDir, "instances")
}

// getDefaultResultsPath returns the default path for results
func getDefaultResultsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, baseDir, "results")
}

// GetHypervisor returns the hypervisor to be used
func getHypervisor() string {
	// Default hypervisor
	hypervisor := "libvirt"

	// If the OS is macOS, use QEMU
	if runtime.GOOS == "darwin" {
		hypervisor = "qemu"
	}

	return hypervisor
}

// GetConnection returns the connection to be used
func getConnection() string {
	// Default connection
	conenction := "qemu:///system"

	// If the OS is macOS, use the connection for QEMU
	if runtime.GOOS == "darwin" {
		conenction = ""
	}

	return conenction
}
