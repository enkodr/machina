package hypvsr

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/imgutil"
	"github.com/enkodr/machina/internal/netutil"
	"github.com/enkodr/machina/internal/osutil"
	"github.com/enkodr/machina/internal/sshutil"
	"github.com/enkodr/machina/internal/usrutil"
	"gopkg.in/yaml.v3"
)

// Machine holds the configuration details for a single machine
type Machine struct {
	baseDir     string
	Kind        string      `yaml:"kind"`                  // Kind of the resource, should be 'Machine'
	Name        string      `yaml:"name,omitempty"`        // Name of the machine. Must be unique in the system
	Extends     string      `yaml:"extends,omitempty"`     // Name of the Machine to extend
	Replicas    int         `yaml:"replicas,omitempty"`    // Number of Replicas (used with kind Cluster)
	Image       Image       `yaml:"image,omitempty"`       // Image details for the machine
	Credentials Credentials `yaml:"credentials,omitempty"` // Credentials for the machine
	Resources   Resources   `yaml:"resources,omitempty"`   // Hardware resources for the machine
	Scripts     Scripts     `yaml:"scripts,omitempty"`     // Scripts to run in the machine
	Mount       Mount       `yaml:"mount,omitempty"`       // Mount point details
	Network     Network     `yaml:"network,omitempty"`     // Network configuration
	Connection  string      `yaml:"connection,omitempty"`  // Connection to hypervisor
	Variant     string      `yaml:"variant,omitempty"`     // OS variant to use
}

// Image holds the URL and checksum of the machine image
type Image struct {
	URL      string `yaml:"url,omitempty"`      // URL of the machine image
	Checksum string `yaml:"checksum,omitempty"` // Checksum for the image in the format 'algorithm:hash'
}

// Credentials holds the username, password, and user groups
type Credentials struct {
	Username string   `yaml:"username,omitempty"` // Username for the machine
	Password string   `yaml:"password,omitempty"` // Password for the machine
	Groups   []string `yaml:"groups,omitempty"`   // User groups for the machine
}

// Resources holds the hardware specifications of the machine
type Resources struct {
	CPUs   string `yaml:"cpus,omitempty"`   // Number of CPUs for the machine
	Memory string `yaml:"memory,omitempty"` // Amount of RAM for the machine
	Disk   string `yaml:"disk,omitempty"`   // Disk space for the machine
}

// Scripts holds the installation and initialisation scripts
type Scripts struct {
	Install string `yaml:"install,omitempty"` // Installation script to run in the machine
	Init    string `yaml:"init,omitempty"`    // Initialisation script to run when machine starts
}

// Mount holds the hostPath and guestPath for mounting host folders into the VM
type Mount struct {
	Name      string `yaml:"name,omitempty"`      // Name of the mount point
	HostPath  string `yaml:"hostPath,omitempty"`  // Path in the host
	GuestPath string `yaml:"guestPath,omitempty"` // Path inside the VM
}

// Network holds the network configuration
type Network struct {
	NicName    string `yaml:"nicName,omitempty"`    // Name of the interface
	IPAddress  string `yaml:"ipAddress,omitempty"`  // IP Address of the machine
	Gateway    string `yaml:"gateway,omitempty"`    // Gateway of the network
	MacAddress string `yaml:"macAddress,omitempty"` // MacAddress of the NIC
}

func (vm *Machine) CreateDir() error {
	// Check if VM already exists
	_, err := os.Stat(filepath.Join(vm.baseDir, vm.Name))
	if !os.IsNotExist(err) {
		return errors.New("machine already exists")
	}

	osutil.MkDir(filepath.Join(vm.baseDir, vm.Name))

	return nil
}

func (vm *Machine) Prepare() error {
	// Credentials network configuration
	net := netutil.NewNetwork()
	netYaml, err := yaml.Marshal(net)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.NetworkFilename)), netYaml, 0644)
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
	err = os.WriteFile(filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.UserdataFilename)), usrYaml, 0644)
	if err != nil {
		return err
	}

	// Save private key
	err = os.WriteFile(filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.PrivateKeyFilename)), clCfg.PrivateKey, 0600)
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

	err = os.WriteFile(filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.InstanceFilename)), vmYaml, 0644)
	if err != nil {
		return err
	}

	return nil

}

func (vm *Machine) DownloadImage() error {
	cfg, _ := config.LoadConfig()
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

func (vm *Machine) CreateDisks() error {
	cfg, _ := config.LoadConfig()
	// Create main disk
	image, _ := imgutil.GetFilenameFromURL(vm.Image.URL)
	command := "qemu-img"
	args := []string{
		"create",
		"-F", "qcow2",
		"-b", filepath.Join(cfg.Directories.Images, image),
		"-f", "qcow2", filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.DiskFilename)),
		vm.Resources.Disk,
	}
	fmt.Println(args)
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// create seed disk
	command = "cloud-localds"
	args = []string{
		fmt.Sprintf("--network-config=%s", filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.NetworkFilename))),
		filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.SeedImageFilename)),
		filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.UserdataFilename)),
	}
	cmd = exec.Command(command, args...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Create creates the VM and starts it
func (vm *Machine) Create() error {
	cfg, _ := config.LoadConfig()
	// Define the Hypervisor to use
	var h Hypervisor
	if cfg.Hypervisor == "qemu" {
		h = &Qemu{}
	} else {
		h = &Libvirt{}
	}
	return h.Create(vm)
}

// Wait until the machine is running
func (vm *Machine) Wait() error {
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
func (vm *Machine) Start() error {
	cfg, _ := config.LoadConfig()
	// Define the Hypervisor to use
	var h Hypervisor
	if cfg.Hypervisor == "qemu" {
		h = &Qemu{}
	} else {
		h = &Libvirt{}
	}
	return h.Start(vm)
}

// Stops a running VM
func (vm *Machine) Stop() error {
	cfg, _ := config.LoadConfig()
	// Define the Hypervisor to use
	var h Hypervisor
	if cfg.Hypervisor == "qemu" {
		h = &Qemu{}
	} else {
		h = &Libvirt{}
	}
	return h.Stop(vm)
}

// Force stops a running VM
func (vm *Machine) ForceStop() error {
	cfg, _ := config.LoadConfig()
	// Define the Hypervisor to use
	var h Hypervisor
	if cfg.Hypervisor == "qemu" {
		h = &Qemu{}
	} else {
		h = &Libvirt{}
	}
	return h.ForceStop(vm)
}

// Gets the status of a VM
func (vm *Machine) Status() (string, error) {
	cfg, _ := config.LoadConfig()
	// Define the Hypervisor to use
	var h Hypervisor
	if cfg.Hypervisor == "qemu" {
		h = &Qemu{}
	} else {
		h = &Libvirt{}
	}
	return h.Status(vm)
}

// Deletes a VM
func (vm *Machine) Delete() error {
	cfg, _ := config.LoadConfig()
	// Define the Hypervisor to use
	var h Hypervisor
	if cfg.Hypervisor == "qemu" {
		h = &Qemu{}
	} else {
		h = &Libvirt{}
	}
	return h.Delete(vm)
}

// Copies content from host to guest or vice-versa
func (vm *Machine) CopyContent(origin string, dest string) error {
	// Define the origin and destination for copying content
	hostToVM := true
	if hostToVM {
		parts := strings.Split(dest, ":")
		dest = fmt.Sprintf("%s@%s:%s", vm.Credentials.Username, vm.Network.IPAddress, parts[1])
	} else {
		parts := strings.Split(origin, ":")
		origin = fmt.Sprintf("%s@%s:%s", vm.Credentials.Username, vm.Network.IPAddress, parts[1])
	}
	command := "scp"
	args := []string{
		"-r",
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.PrivateKeyFilename)),
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

// Runs the initial scripts after the machine is created
func (vm *Machine) RunInitScripts() error {
	// Copies the scripts
	command := "scp"
	args := []string{
		"-o", "StrictHostKeyChecking=no",
		"-i", filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.PrivateKeyFilename)),
		"-r",
		filepath.Join(vm.baseDir, vm.Name, "bin/"),
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
		"-i", filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.PrivateKeyFilename)),
		fmt.Sprintf("%s@%s", vm.Credentials.Username, vm.Network.IPAddress),
		"/tmp/machina/install.sh",
	}
	cmd = exec.Command(command, args...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Cleanup
	return os.RemoveAll(filepath.Join(vm.baseDir, vm.Name, "bin"))
}

func (vm *Machine) Shell() error {
	// Copies the init script into the VM
	command := "ssh"
	args := []string{
		"-i", filepath.Join(vm.baseDir, vm.Name, config.GetFilename(config.PrivateKeyFilename)),
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

func (vm *Machine) GetVMs() []Machine {
	return []Machine{*vm}
}

// Creates the script files from the template
func (vm *Machine) createScriptFiles() error {
	err := os.MkdirAll(filepath.Join(vm.baseDir, vm.Name, "bin"), 0755)
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

	err = os.WriteFile(filepath.Join(vm.baseDir, vm.Name, "bin/machina.service"), []byte(sysDSvc), 0744)
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
	err = os.WriteFile(filepath.Join(vm.baseDir, vm.Name, "bin/install.sh"), []byte(vm.Scripts.Install), 0744)
	if err != nil {
		return nil
	}

	// Boot script
	var mountName string
	cfg, err := config.LoadConfig()
	if cfg != nil {
		return err
	}
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

	err = os.WriteFile(filepath.Join(vm.baseDir, vm.Name, "bin/machina"), []byte(initScript), 0744)
	if err != nil {
		return nil
	}

	err = os.WriteFile(filepath.Join(vm.baseDir, vm.Name, "bin/machinarc"), []byte(vm.Scripts.Init), 0644)
	if err != nil {
		return nil
	}
	return nil
}
