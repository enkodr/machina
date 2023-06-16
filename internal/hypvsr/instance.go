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

// Instance holds the configuration details for a single machine
type Instance struct {
	baseDir     string        `yaml:"-"`
	Kind        string        `yaml:"kind"`                  // Kind of the resource, should be 'Machine'
	Name        string        `yaml:"name,omitempty"`        // Name of the machine. Must be unique in the system
	Extends     string        `yaml:"extends,omitempty"`     // Name of the Machine to extend
	Replicas    int           `yaml:"replicas,omitempty"`    // Number of Replicas (used with kind Cluster)
	Image       Image         `yaml:"image,omitempty"`       // Image details for the machine
	Credentials Credentials   `yaml:"credentials,omitempty"` // Credentials for the machine
	Resources   Resources     `yaml:"resources,omitempty"`   // Hardware resources for the machine
	Scripts     Scripts       `yaml:"scripts,omitempty"`     // Scripts to run in the machine
	Mount       Mount         `yaml:"mount,omitempty"`       // Mount point details
	Network     Network       `yaml:"network,omitempty"`     // Network configuration
	Connection  string        `yaml:"connection,omitempty"`  // Connection to hypervisor
	Variant     string        `yaml:"variant,omitempty"`     // OS variant to use
	Hypervisor  Hypervisor    `yaml:"-"`
	Runner      osutil.Runner `yaml:"-"`
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

// CreateDir creates the directory for the machine
func (instance *Instance) CreateDir() error {
	// Check if VM already exists
	_, err := os.Stat(filepath.Join(instance.baseDir, instance.Name))
	if !os.IsNotExist(err) {
		return errors.New("machine already exists")
	}

	// Create the directory
	return os.Mkdir(filepath.Join(instance.baseDir, instance.Name), 0755)

}

// Prepare prepares the machine for use
func (instance *Instance) Prepare() error {
	// Credentials network configuration
	net := netutil.NewNetwork()
	netYaml, err := yaml.Marshal(net)
	if err != nil {
		return err
	}

	// Save network configuration
	err = os.WriteFile(filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.NetworkFilename)), netYaml, 0644)
	if err != nil {
		return err
	}

	// Get the IP address from the network configuration
	ipAddr, err := netutil.GetIPFromNetworkAddress(net.Ethernets.VirtNet.Addresses[0])
	if err != nil {
		return err
	}

	// Set Network configuration
	instance.Network = Network{
		NicName:    net.Ethernets.VirtNet.Name,
		IPAddress:  ipAddr,
		Gateway:    net.Ethernets.VirtNet.Gateway4,
		MacAddress: net.Ethernets.VirtNet.Match.MacAddress,
	}

	// Create user data
	clCfg := usrutil.CloudConfig{
		Hostname: instance.Name,
		Username: instance.Credentials.Username,
		Password: instance.Credentials.Password,
		Groups:   instance.Credentials.Groups,
	}

	// Create user data
	usr, err := usrutil.NewUserData(&clCfg)
	if err != nil {
		return err
	}

	// Save user data
	usrYaml, err := yaml.Marshal(usr)
	if err != nil {
		return err
	}

	usrYaml = append([]byte("#cloud-config\n"), usrYaml...)
	err = os.WriteFile(filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.UserdataFilename)), usrYaml, 0644)
	if err != nil {
		return err
	}

	// Save private key
	err = os.WriteFile(filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.PrivateKeyFilename)), clCfg.PrivateKey, 0600)
	if err != nil {
		return err
	}

	// Create script files
	err = instance.createScriptFiles()
	if err != nil {
		return err
	}

	// Save machine file
	vmYaml, err := yaml.Marshal(instance)
	if err != nil {
		return err
	}

	// Set the hypervisor
	instance.Hypervisor = getHypervisor()

	err = os.WriteFile(filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.InstanceFilename)), vmYaml, 0644)
	if err != nil {
		return err
	}

	return nil

}

// DownloadImage downloads the image for the machine
func (instance *Instance) DownloadImage() error {
	// Get the image filename
	imgDir := cfg.Directories.Images
	fileName, err := imgutil.GetFilenameFromURL(instance.Image.URL)
	if err != nil {
		return err
	}

	// Set the local image path
	localImage := filepath.Join(imgDir, fileName)

	// check if hashes equal
	if osutil.Checksum(localImage, instance.Image.Checksum) {
		return nil
	}

	// download the image
	err = netutil.DownloadAndSave(instance.Image.URL, imgDir)
	if err != nil {
		return err
	}
	return nil
}

// CreateDisks creates the disks for the machine
func (instance *Instance) CreateDisks() error {
	err := instance.createInstanceDisk()
	if err != nil {
		return err
	}

	err = instance.createSeedDisk()
	if err != nil {
		return err
	}

	return nil
}

func (instance *Instance) createInstanceDisk() error {
	// Get the image filename
	image, _ := imgutil.GetFilenameFromURL(instance.Image.URL)

	// Set the command
	command := "qemu-img"

	// Set the arguments
	args := []string{
		"create",
		"-F", "qcow2",
		"-b", filepath.Join(cfg.Directories.Images, image),
		"-f", "qcow2", filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.DiskFilename)),
		instance.Resources.Disk,
	}

	// Run the command to create the disk
	_, err := instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return nil
}

func (instance *Instance) createSeedDisk() error {
	// Set the command
	command := "cloud-localds"

	// Set the arguments
	args := []string{
		fmt.Sprintf("--network-config=%s", filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.NetworkFilename))),
		filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.SeedImageFilename)),
		filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.UserdataFilename)),
	}

	// Run the command to create the disk
	_, err := instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return nil
}

// Create creates the VM and starts it
func (instance *Instance) Create() error {
	return instance.Hypervisor.Create(instance)
}

// Wait until the machine is running
func (instance *Instance) Wait() error {
	// Set the start time
	start := time.Now()
	running := false
	for !running {
		// Check if the machine is running
		running = sshutil.IsResponding(instance.Network.IPAddress)
		// Sleep for 1 second
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
func (instance *Instance) Start() error {
	return instance.Hypervisor.Start(instance)
}

// Stops a running VM
func (instance *Instance) Stop() error {
	return instance.Hypervisor.Stop(instance)
}

// Force stops a running VM
func (instance *Instance) ForceStop() error {
	return instance.Hypervisor.ForceStop(instance)
}

// Gets the status of a VM
func (instance *Instance) Status() (string, error) {
	return instance.Hypervisor.Status(instance)
}

// Deletes a VM
func (instance *Instance) Delete() error {
	return instance.Hypervisor.Delete(instance)
}

// Copies content from host to guest or vice-versa
func (instance *Instance) CopyContent(origin string, dest string) error {
	// Define the origin and destination for copying content
	hostToVM := true
	if hostToVM {
		parts := strings.Split(dest, ":")
		dest = fmt.Sprintf("%s@%s:%s", instance.Credentials.Username, instance.Network.IPAddress, parts[1])
	} else {
		parts := strings.Split(origin, ":")
		origin = fmt.Sprintf("%s@%s:%s", instance.Credentials.Username, instance.Network.IPAddress, parts[1])
	}
	command := "scp"
	args := []string{
		"-r",
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.PrivateKeyFilename)),
		origin,
		dest,
	}

	// Run the command to create the disk
	_, err := instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return nil
}

// Runs the initial scripts after the machine is created
func (instance *Instance) RunInitScripts() error {
	// Copies the scripts
	command := "scp"
	args := []string{
		"-o", "StrictHostKeyChecking=no",
		"-i", filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.PrivateKeyFilename)),
		"-r",
		filepath.Join(instance.baseDir, instance.Name, "bin/"),
		fmt.Sprintf("%s@%s:/tmp/machina", instance.Credentials.Username, instance.Network.IPAddress),
	}

	// Run the command to create the disk
	_, err := instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Runs the init script inside the VM
	command = "ssh"
	args = []string{
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.PrivateKeyFilename)),
		fmt.Sprintf("%s@%s", instance.Credentials.Username, instance.Network.IPAddress),
		"/tmp/machina/install.sh",
	}

	// Run the command to create the disk
	_, err = instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Cleanup
	return os.RemoveAll(filepath.Join(instance.baseDir, instance.Name, "bin"))
}

func (instance *Instance) Shell() error {
	// Copies the init script into the VM
	command := "ssh"
	args := []string{
		"-i", filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.PrivateKeyFilename)),
		fmt.Sprintf("%s@%s", instance.Credentials.Username, instance.Network.IPAddress),
	}

	// TODO: Add stdin, stdout, stderr to os.Stdout, os.Stdin, os.Stderr
	// Run the command to create the disk
	// _, err := instance.Runner.RunCommand(command, args, osutil.WithStdin(os.Stdin), osutil.WithStdout(os.Stdout), osutil.WithStderr(os.Stderr))
	// if err != nil {
	// 	return err
	// }
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

func (instance *Instance) GetVMs() []Instance {
	return []Instance{*instance}
}

// Creates the script files from the template
func (instance *Instance) createScriptFiles() error {
	err := os.MkdirAll(filepath.Join(instance.baseDir, instance.Name, "bin"), 0755)
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

	err = os.WriteFile(filepath.Join(instance.baseDir, instance.Name, "bin/machina.service"), []byte(sysDSvc), 0744)
	if err != nil {
		return nil
	}

	// Install script
	installScript := `
sudo mkdir -p /etc/machina
echo 'source /etc/machina/machinarc' >> $HOME/.bashrc
sudo mv /tmp/machina/* /etc/machina
sudo chcon -R -t bin_t /etc/machina/machina.service
sudo cp /etc/machina/machina.service /etc/systemd/system/machina.service
sudo chmod 664 /etc/systemd/system/machina.service
sudo systemctl daemon-reload
sudo systemctl enable machina.service
sudo systemctl start machina.service
`
	instance.Scripts.Install += installScript
	err = os.WriteFile(filepath.Join(instance.baseDir, instance.Name, "bin/install.sh"), []byte(instance.Scripts.Install), 0744)
	if err != nil {
		return nil
	}

	// Boot script
	var mountName string
	if cfg.Hypervisor == "qemu" {
		mountName = instance.Mount.Name
	} else {
		mountName = instance.Mount.GuestPath
	}
	initScript := fmt.Sprintf(`#!/bin/bash
HOST_PATH="%s"
GUEST_PATH="%s"
if [[ "$GUEST_PATH" != "" ]]; then
	mkdir -p $GUEST_PATH
	sudo mount -t 9p $HOST_PATH $GUEST_PATH
fi

exit 0
`, mountName, instance.Mount.GuestPath)

	err = os.WriteFile(filepath.Join(instance.baseDir, instance.Name, "bin/machina"), []byte(initScript), 0744)
	if err != nil {
		return nil
	}

	err = os.WriteFile(filepath.Join(instance.baseDir, instance.Name, "bin/machinarc"), []byte(instance.Scripts.Init), 0644)
	if err != nil {
		return nil
	}
	return nil
}
