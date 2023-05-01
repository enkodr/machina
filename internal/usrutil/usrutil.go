package usrutil

import (
	"fmt"
	"strings"

	"github.com/enkodr/machina/internal/sshutil"
)

type CloudConfig struct {
	Hostname   string
	Username   string
	Password   string
	Groups     []string
	PrivateKey []byte
}
type UserData struct {
	Hostname       string     `yaml:"hostname"`
	ManageEtcHosts bool       `yaml:"manage_etc_hosts"`
	Users          []User     `yaml:"users"`
	SSHPwAuth      bool       `yaml:"ssh_pwauth"`
	DisableRoot    bool       `yaml:"disable_root"`
	ChPassword     ChPassword `yaml:"chpasswd"`
}

type User struct {
	Name              string   `yaml:"name"`
	Password          string   `yaml:"password"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
	Sudo              string   `yaml:"sudo"`
	Groups            string   `yaml:"groups"`
	Home              string   `yaml:"home"`
	Shell             string   `yaml:"shell"`
	LockPassword      bool     `yaml:"lock_passwd"`
}

type ChPassword struct {
	List   string `yaml:"list"`
	Expire bool   `yaml:"expire"`
}

func NewUserData(cfg *CloudConfig) (*UserData, error) {
	priv, pub, err := sshutil.GenerateNewSSHKeys()
	if err != nil {
		return nil, err
	}

	cfg.PrivateKey = priv

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
