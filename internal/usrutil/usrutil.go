package usrutil

import (
	"fmt"
	"strings"

	"github.com/enkodr/machina/internal/sshutil"
)

// CloudConfig is a struct that holds the configuration for the cloud-config file
type CloudConfig struct {
	Hostname   string   // Hostname is the hostname of the instance
	Username   string   // Username is the username of the instance
	Password   string   // Password is the password of the instance
	Groups     []string //	Groups is the groups of the instance
	PrivateKey []byte   // PrivateKey is the private key of the instance
}

// UserData is a struct that holds the configuration for the user-data file
type UserData struct {
	Hostname       string     `yaml:"hostname"`         // Hostname is the hostname of the userdata
	ManageEtcHosts bool       `yaml:"manage_etc_hosts"` // ManageEtcHosts is a boolean that determines if the /etc/hosts file should be managed
	Users          []User     `yaml:"users"`            // Users is a slice of users
	SSHPwAuth      bool       `yaml:"ssh_pwauth"`       // SSHPwAuth is a boolean that determines if password authentication is allowed
	DisableRoot    bool       `yaml:"disable_root"`     // DisableRoot is a boolean that determines if the root user should be disabled
	ChPassword     ChPassword `yaml:"chpasswd"`         // ChPassword is a struct that holds the configuration for the chpasswd module
}

// User is a struct that holds the configuration for the user-data file
type User struct {
	Name              string   `yaml:"name"`                // Name is the name of the user
	Password          string   `yaml:"password"`            // Password is the password of the user
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"` // SSHAuthorizedKeys is a slice of ssh public keys
	Sudo              string   `yaml:"sudo"`                // Sudo is the sudo configuration for the user
	Groups            string   `yaml:"groups"`              // Groups is the groups of the user
	Home              string   `yaml:"home"`                // Home is the home directory of the user
	Shell             string   `yaml:"shell"`               // Shell is the shell of the user
	LockPassword      bool     `yaml:"lock_passwd"`         // LockPassword is a boolean that determines if the password should be locked
}

// ChPassword is a struct that holds the configuration for the chpasswd module
type ChPassword struct {
	List   string `yaml:"list"`   // List is a list of users and passwords
	Expire bool   `yaml:"expire"` // Expire is a boolean that determines if the password should be expired
}

// NewUserData creates a new UserData struct
func NewUserData(cfg *CloudConfig) (*UserData, error) {
	// Generate a new SSH key pair
	priv, pub, err := sshutil.GenerateNewSSHKeys()
	if err != nil {
		return nil, err
	}

	// Set the private key in the cloud config
	cfg.PrivateKey = priv

	// Create a new UserData struct
	usr := &UserData{
		Hostname:       cfg.Hostname,
		ManageEtcHosts: true,
		Users: []User{
			{
				Name:     cfg.Username,
				Password: cfg.Password,
				SSHAuthorizedKeys: []string{
					string(pub),
				},
				Sudo:         "ALL=(ALL) NOPASSWD:ALL",
				Groups:       strings.Join(cfg.Groups, ", "),
				Home:         fmt.Sprintf("/home/%s", cfg.Username),
				Shell:        "/bin/bash",
				LockPassword: false,
			},
		},
		SSHPwAuth:   true,
		DisableRoot: false,
		ChPassword: ChPassword{
			List:   fmt.Sprintf("%s:%s", cfg.Username, cfg.Username),
			Expire: false,
		},
	}

	return usr, nil
}
