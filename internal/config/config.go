package config

import (
	"encoding/json"
	"runtime"

	"github.com/enkodr/machina/internal/db"
	"github.com/enkodr/machina/internal/path"
)

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

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	// If the config file exists, load it
	cfg, err := loadConfigFromDB()
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		return createDefaultConfig()
	}

	return cfg, nil
}

// loadConfig loads the configuration from the config file
func loadConfigFromDB() (*Config, error) {
	// Initialise db
	db, err := db.NewDB("config")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	data, err := db.Get("config")
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	// Unmarshal the config file
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// createDefaultConfig creates a default config file
func createDefaultConfig() (*Config, error) {

	cfg := Config{
		Hypervisor: getHypervisor(),
		Connection: getConnection(),
		Directories: Directories{
			Images:    path.GetPath(path.ImagesDir),
			Instances: path.GetPath(path.InstanceFile),
			Results:   path.GetPath(path.ResultsDir),
		},
	}

	db, err := db.NewDB("config")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Create the config file
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	db.Put("config", []byte(data))

	return &cfg, nil
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
