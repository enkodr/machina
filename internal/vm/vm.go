package vm

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/imgutil"
	"github.com/enkodr/machina/internal/netutil"
	"github.com/enkodr/machina/internal/osutil"
	"github.com/enkodr/machina/internal/sshutil"
	"github.com/enkodr/machina/internal/usrutil"
	"gopkg.in/yaml.v3"
)

// VMConfig represents a machina virtual machine
type VMConfig struct {
	Hypervisor  Hypervisor  `yaml:"-"`
	Extends     string      `yaml:"extends,omitempty"`
	Name        string      `yaml:"name,omitempty"`
	Image       Image       `yaml:"image,omitempty"`
	Credentials Credentials `yaml:"credentials,omitempty"`
	Specs       Specs       `yaml:"specs,omitempty"`
	Scripts     Scripts     `yaml:"scripts,omitempty"`
	Mount       Mount       `yaml:"mount,omitempty"`
	Network     Network     `yaml:"network,omitempty"`
	Connection  string      `yaml:"connection,omitempty"`
	Variant     string      `yaml:"variant,omitempty"`
}

// Image represents the distro image to use
type Image struct {
	URL      string `yaml:"url,omitempty"`
	Checksum string `yaml:"checksum,omitempty"`
}

// Credentials represents the credentials for the machine user
type Credentials struct {
	Username string   `yaml:"username,omitempty"`
	Password string   `yaml:"password,omitempty"`
	Groups   []string `yaml:"groups,omitempty"`
}

// Specs represents the hardware specifications
type Specs struct {
	CPUs   string `yaml:"cpus,omitempty"`
	Memory string `yaml:"memory,omitempty"`
	Disk   string `yaml:"disk,omitempty"`
}

// Scripts represents the scripts to run in the machine
type Scripts struct {
	Install string `yaml:"install,omitempty"`
	Init    string `yaml:"init,omitempty"`
}

// Mount represents the mount points from the to the machine
type Mount struct {
	Name      string `yaml:"name,omitempty"`
	HostPath  string `yaml:"hostPath,omitempty"`
	GuestPath string `yaml:"guestPath,omitempty"`
}

// Network represents the network configuration
type Network struct {
	NicName    string `yaml:"nicName,omitempty"`
	IPAddress  string `yaml:"ipAddress,omitempty"`
	Gateway    string `yaml:"gateway,omitempty"`
	MacAddress string `yaml:"macAddress,omitempty"`
}

type Hypervisor interface {
	Create(vm *VMConfig) error
	Start(vm *VMConfig) error
	Stop(vm *VMConfig) error
	ForceStop(vm *VMConfig) error
	Status(vm *VMConfig) (string, error)
	Delete(vm *VMConfig) error
}

func NewVM(name string) (*VMConfig, error) {
	tpl := NewTemplate(name)
	vm, err := tpl.Load()
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func (vm *VMConfig) Prepare() error {
	// Load configuration
	cfg := config.LoadConfig()

	// Check if VM already exists
	_, err := os.Stat(filepath.Join(cfg.Directories.Instances, vm.Name))
	if !os.IsNotExist(err) {
		return errors.New("machine already exists")
	}

	if cfg.Hypervisor == "qemu" {
		vm.Hypervisor = &Qemu{}
	} else {
		vm.Hypervisor = &Libvirt{}
	}

	// Create directory structure
	imgutil.EnsureDirectories(vm.Name)

	// Create network configuration
	net := netutil.NewNetwork()
	netYaml, err := yaml.Marshal(net)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(cfg.Directories.Instances, vm.Name, "network.cfg"), netYaml, 0644)
	if err != nil {
		return err
	}
	vm.Network = Network{
		NicName:    net.Ethernets.VirtNet.Name,
		IPAddress:  netutil.GetIPFromNetworkAddress(net.Ethernets.VirtNet.Addresses[0]),
		Gateway:    net.Ethernets.VirtNet.Gateway4,
		MacAddress: net.Ethernets.VirtNet.Match.MacAddress,
	}

	// Create user data
	clCfg := usrutil.CloudConfig{
		Hostname: vm.Name,
		Username: vm.Credentials.Username,
		Password: vm.Credentials.Password,
		Groups:   vm.Credentials.Groups,
	}
	usr, err := usrutil.NewUserData(&clCfg)
	if err != nil {
		return err
	}
	usrYaml, err := yaml.Marshal(usr)
	if err != nil {
		return err
	}
	usrYaml = append([]byte("#cloud-config\n"), usrYaml...)
	err = os.WriteFile(filepath.Join(cfg.Directories.Instances, vm.Name, "userdata.yaml"), usrYaml, 0644)
	if err != nil {
		return err
	}

	// Save private key
	err = os.WriteFile(filepath.Join(cfg.Directories.Instances, vm.Name, "id_rsa"), clCfg.PrivateKey, 0600)
	if err != nil {
		return err
	}

	// Create script files
	vm.createScriptFiles()

	// Save machine file
	vmYaml, err := yaml.Marshal(vm)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(cfg.Directories.Instances, vm.Name, "machina.yaml"), vmYaml, 0644)
	if err != nil {
		return err
	}

	return nil

}

func Load(name string) (*VMConfig, error) {
	cfg := config.LoadConfig()
	vm := &VMConfig{}

	data, err := os.ReadFile(filepath.Join(cfg.Directories.Instances, name, "machina.yaml"))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, vm)
	if err != nil {
		return nil, err
	}
	if cfg.Hypervisor == "qemu" {
		vm.Hypervisor = &Qemu{}
	} else {
		vm.Hypervisor = &Libvirt{}
	}
	return vm, nil
}

func (vm *VMConfig) DownloadImage() error {
	cfg := config.LoadConfig()
	imgDir := cfg.Directories.Images
	fileName, err := imgutil.GetFilenameFromURL(vm.Image.URL)
	if err != nil {
		return err
	}

	localImage := filepath.Join(imgDir, fileName)

	// check if hashes equal
	if osutil.Checksum(localImage, vm.Image.Checksum) {
		return nil
	}

	// download the image
	err = netutil.DownloadAndSave(vm.Image.URL, imgDir)
	if err != nil {
		return err
	}
	return nil
}

func (vm *VMConfig) CreateDisks() error {
	// Create main disk
	image, _ := imgutil.GetFilenameFromURL(vm.Image.URL)
	cfg := config.LoadConfig()
	command := "qemu-img"
	args := []string{
		"create",
		"-F", "qcow2",
		"-b", filepath.Join(cfg.Directories.Images, image),
		"-f", "qcow2", filepath.Join(cfg.Directories.Instances, vm.Name, "disk.img"),
		vm.Specs.Disk,
	}
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// create seed disk
	command = "cloud-localds"
	args = []string{
		fmt.Sprintf("--network-config=%s", filepath.Join(cfg.Directories.Instances, vm.Name, "network.cfg")),
		filepath.Join(cfg.Directories.Instances, vm.Name, "seed.img"),
		filepath.Join(cfg.Directories.Instances, vm.Name, "userdata.yaml"),
	}
	cmd = exec.Command(command, args...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Create creates the VM and starts it
func (vm *VMConfig) Create() error {
	return vm.Hypervisor.Create(vm)
}

// Wait until the machine is running
func (vm *VMConfig) Wait() error {
	start := time.Now()
	running := false
	for !running {
		running = sshutil.IsResponding(vm.Network.IPAddress)
		time.Sleep(time.Second)
		// Return a timeout error in case the machine takes more than
		// 5 minutes to become responsive
		if time.Since(start) >= time.Second*300 {
			return errors.New("timeout")
		}
	}

	return nil
}

// Starts a stopped vm
func (vm *VMConfig) Start() error {
	return vm.Hypervisor.Start(vm)
}

// Stops a running VM
func (vm *VMConfig) Stop() error {
	return vm.Hypervisor.Stop(vm)
}

// Force stops a running VM
func (vm *VMConfig) ForceStop() error {
	return vm.Hypervisor.ForceStop(vm)
}

// Gets the status of a VM
func (vm *VMConfig) Status() (string, error) {
	return vm.Hypervisor.Status(vm)
}

// Deletes a VM
func (vm *VMConfig) Delete() error {
	return vm.Hypervisor.Delete(vm)
}

// Copies content from host to guest or vice-versa
func (vm *VMConfig) CopyContent(origin string, dest string) error {
	cfg := config.LoadConfig()
	command := "scp"
	args := []string{
		"-r",
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(cfg.Directories.Instances, vm.Name, "id_rsa"),
		origin,
		dest,
	}

	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Creates the script files from the template
func (vm *VMConfig) createScriptFiles() error {
	// Prepare
	cfg := config.LoadConfig()
	err := os.MkdirAll(filepath.Join(cfg.Directories.Instances, vm.Name, "bin"), 0755)
	if err != nil {
		return nil
	}

	// Systemd service
	// sudo journalctl -xeu machina.service
	sysDSvc := `[Unit]
Description=machina mount

[Service]
Type=forking
User=machina
Group=machina
ExecStart=/etc/machina/machina
StandardOutput=journal
	
[Install]
WantedBy=multi-user.target
`

	err = os.WriteFile(filepath.Join(cfg.Directories.Instances, vm.Name, "bin/machina.service"), []byte(sysDSvc), 0744)
	if err != nil {
		return nil
	}

	// Install script
	installScript := `
sudo mkdir -p /etc/machina
echo 'source /etc/machina/machinarc' >> $HOME/.bashrc
sudo mv /tmp/machina/* /etc/machina
sudo chcon -R -t bin_t /etc/machina/machina
sudo cp /etc/machina/machina.service /etc/systemd/system/machina.service
sudo chmod 664 /etc/systemd/system/machina.service
sudo systemctl daemon-reload
sudo systemctl enable machina.service
sudo systemctl start machina.service
`
	vm.Scripts.Install += installScript
	err = os.WriteFile(filepath.Join(cfg.Directories.Instances, vm.Name, "bin/install.sh"), []byte(vm.Scripts.Install), 0744)
	if err != nil {
		return nil
	}

	// Boot script
	var mountName string
	if cfg.Hypervisor == "qemu" {
		mountName = vm.Mount.Name
	} else {
		mountName = vm.Mount.GuestPath
	}
	initScript := fmt.Sprintf(`#!/bin/bash
HOST_PATH="%s"
GUEST_PATH="%s"
if [[ "$GUEST_PATH" != "" ]]; then
	mkdir -p $GUEST_PATH
	sudo mount -t 9p $HOST_PATH $GUEST_PATH
fi

exit 0
`, mountName, vm.Mount.GuestPath)

	err = os.WriteFile(filepath.Join(cfg.Directories.Instances, vm.Name, "bin/machina"), []byte(initScript), 0744)
	if err != nil {
		return nil
	}

	err = os.WriteFile(filepath.Join(cfg.Directories.Instances, vm.Name, "bin/machinarc"), []byte(vm.Scripts.Init), 0644)
	if err != nil {
		return nil
	}
	return nil
}

// Runs the initial scripts after the machine is created
func (vm *VMConfig) RunInitScripts() error {
	cfg := config.LoadConfig()

	// Copies the scripts
	command := "scp"
	args := []string{
		"-o", "StrictHostKeyChecking=no",
		"-i", filepath.Join(cfg.Directories.Instances, vm.Name, "id_rsa"),
		"-r",
		filepath.Join(cfg.Directories.Instances, vm.Name, "bin/"),
		fmt.Sprintf("%s@%s:/tmp/machina", vm.Credentials.Username, vm.Network.IPAddress),
	}

	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Runs the init script inside the VM
	command = "ssh"
	args = []string{
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(cfg.Directories.Instances, vm.Name, "id_rsa"),
		fmt.Sprintf("%s@%s", vm.Credentials.Username, vm.Network.IPAddress),
		"/tmp/machina/install.sh",
	}
	cmd = exec.Command(command, args...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Cleanup
	return os.RemoveAll(filepath.Join(cfg.Directories.Instances, vm.Name, "bin"))
}

func (vm *VMConfig) Shell() error {
	cfg := config.LoadConfig()

	// Copies the init script into the VM
	command := "ssh"
	args := []string{
		"-i", filepath.Join(cfg.Directories.Instances, vm.Name, "id_rsa"),
		fmt.Sprintf("%s@%s", vm.Credentials.Username, vm.Network.IPAddress),
	}
	cmd := exec.Command(command, args...)
	// Redirect stdin, stdout and stderr from the SSH connection to the host
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
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
