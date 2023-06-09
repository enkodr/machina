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
	clusterName string
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
func (machine *Machine) CreateDir() error {
	// Check if VM already exists
	_, err := os.Stat(filepath.Join(machine.baseDir, machine.Name))
	if !os.IsNotExist(err) {
		return errors.New("machine already exists")
	}

	// Create the directory
	return os.Mkdir(filepath.Join(machine.baseDir, machine.Name), 0755)

}

// Prepare prepares the machine for use
func (machine *Machine) Prepare() error {
	// Credentials network configuration
	net := netutil.NewNetwork()
	netYaml, err := yaml.Marshal(net)
	if err != nil {
		return err
	}

	// Save network configuration
	err = os.WriteFile(filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.NetworkFilename)), netYaml, 0644)
	if err != nil {
		return err
	}

	// Get the IP address from the network configuration
	ipAddr, err := netutil.GetIPFromNetworkAddress(net.Ethernets.VirtNet.Addresses[0])
	if err != nil {
		return err
	}

	// Set Network configuration
	machine.Network = Network{
		NicName:    net.Ethernets.VirtNet.Name,
		IPAddress:  ipAddr,
		Gateway:    net.Ethernets.VirtNet.Gateway4,
		MacAddress: net.Ethernets.VirtNet.Match.MacAddress,
	}

	// Create user data
	clCfg := usrutil.CloudConfig{
		Hostname: machine.Name,
		Username: machine.Credentials.Username,
		Password: machine.Credentials.Password,
		Groups:   machine.Credentials.Groups,
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
	err = os.WriteFile(filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.UserdataFilename)), usrYaml, 0644)
	if err != nil {
		return err
	}

	// Save private key
	err = os.WriteFile(filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.PrivateKeyFilename)), clCfg.PrivateKey, 0600)
	if err != nil {
		return err
	}

	// Create script files
	err = machine.createScriptFiles()
	if err != nil {
		return err
	}

	// Save machine file
	vmYaml, err := yaml.Marshal(machine)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.InstanceFilename)), vmYaml, 0644)
	if err != nil {
		return err
	}

	return nil

}

// DownloadImage downloads the image for the machine
func (machine *Machine) DownloadImage() error {
	// Get the image filename
	imgDir := cfg.Directories.Images
	fileName, err := imgutil.GetFilenameFromURL(machine.Image.URL)
	if err != nil {
		return err
	}

	// Set the local image path
	localImage := filepath.Join(imgDir, fileName)

	// check if hashes equal
	if osutil.Checksum(localImage, machine.Image.Checksum) {
		return nil
	}

	// download the image
	err = netutil.DownloadAndSave(machine.Image.URL, imgDir)
	if err != nil {
		return err
	}
	return nil
}

// CreateDisks creates the disks for the machine
func (machine *Machine) CreateDisks() error {
	err := machine.createInstanceDisk()
	if err != nil {
		return err
	}

	err = machine.createSeedDisk()
	if err != nil {
		return err
	}

	return nil
}

func (machine *Machine) createInstanceDisk() error {
	// Get the image filename
	image, _ := imgutil.GetFilenameFromURL(machine.Image.URL)

	// Set the command
	command := "qemu-img"

	// Set the arguments
	args := []string{
		"create",
		"-F", "qcow2",
		"-b", filepath.Join(cfg.Directories.Images, image),
		"-f", "qcow2", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.DiskFilename)),
		machine.Resources.Disk,
	}

	// Run the command to create the disk
	_, err := machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return nil
}

func (machine *Machine) createSeedDisk() error {
	// Set the command
	command := "cloud-localds"

	// Set the arguments
	args := []string{
		fmt.Sprintf("--network-config=%s", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.NetworkFilename))),
		filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.SeedImageFilename)),
		filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.UserdataFilename)),
	}

	// Run the command to create the disk
	_, err := machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return nil
}

// Create creates the VM and starts it
func (machine *Machine) Create() error {
	return machine.Hypervisor.Create(machine)
}

// Wait until the machine is running
func (machine *Machine) Wait() error {
	// Set the start time
	start := time.Now()
	running := false
	for !running {
		// Check if the machine is running
		running = sshutil.IsResponding(machine.Network.IPAddress)
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
func (machine *Machine) Start() error {
	return machine.Hypervisor.Start(machine)
}

// Stops a running VM
func (machine *Machine) Stop() error {
	return machine.Hypervisor.Stop(machine)
}

// Force stops a running VM
func (machine *Machine) ForceStop() error {
	return machine.Hypervisor.ForceStop(machine)
}

// Gets the status of a VM
func (machine *Machine) Status() (string, error) {
	return machine.Hypervisor.Status(machine)
}

// Deletes a VM
func (machine *Machine) Delete() error {
	return machine.Hypervisor.Delete(machine)
}

// Copies content from host to guest or vice-versa
func (machine *Machine) CopyContent(origin string, dest string) error {
	// Define the origin and destination for copying content
	hostToVM := true
	if hostToVM {
		parts := strings.Split(dest, ":")
		dest = fmt.Sprintf("%s@%s:%s", machine.Credentials.Username, machine.Network.IPAddress, parts[1])
	} else {
		parts := strings.Split(origin, ":")
		origin = fmt.Sprintf("%s@%s:%s", machine.Credentials.Username, machine.Network.IPAddress, parts[1])
	}
	command := "scp"
	args := []string{
		"-r",
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.PrivateKeyFilename)),
		origin,
		dest,
	}

	// Run the command to create the disk
	_, err := machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return nil
}

// Runs the initial scripts after the machine is created
func (machine *Machine) RunInitScripts() error {
	// Copies the scripts
	command := "scp"
	args := []string{
		"-o", "StrictHostKeyChecking=no",
		"-i", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.PrivateKeyFilename)),
		"-r",
		filepath.Join(machine.baseDir, machine.Name, "bin/"),
		fmt.Sprintf("%s@%s:/tmp/machina", machine.Credentials.Username, machine.Network.IPAddress),
	}

	// Copy the scripts
	_, err := machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Runs the prepare script inside the VM
	command = "ssh"
	args = []string{
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.PrivateKeyFilename)),
		fmt.Sprintf("%s@%s", machine.Credentials.Username, machine.Network.IPAddress),
		"/tmp/machina/prepare.sh",
	}

	// Run the prepare script
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Set permissions
	command = "chmod"
	args = []string{
		"-R",
		"+x",
		filepath.Join(cfg.Directories.Results, machine.clusterName),
	}

	// Run the set permissions command
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Sync the results from the host to the VM
	command = "rsync"
	args = []string{
		"-ru",
		"-e",
		fmt.Sprintf("ssh -i %s", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.PrivateKeyFilename))),
		filepath.Join(cfg.Directories.Results, machine.clusterName),
		fmt.Sprintf("%s@%s:/etc/machina/results", machine.Credentials.Username, machine.Network.IPAddress),
	}

	// Copy the results to the machine
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Runs the install script inside the Machine
	command = "ssh"
	args = []string{
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.PrivateKeyFilename)),
		fmt.Sprintf("%s@%s", machine.Credentials.Username, machine.Network.IPAddress),
		"/etc/machina/install.sh",
	}

	// Run the install script
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Sync the results from the Machine to the host
	command = "rsync"
	args = []string{
		"-ru",
		"-e",
		fmt.Sprintf("ssh -i %s", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.PrivateKeyFilename))),
		fmt.Sprintf("%s@%s:/etc/machina/results/%s", machine.Credentials.Username, machine.Network.IPAddress, machine.clusterName),
		cfg.Directories.Results,
	}

	// Copy the results to the host
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Cleanup
	return os.RemoveAll(filepath.Join(machine.baseDir, machine.Name, "bin"))
}

func (machine *Machine) Shell() error {
	// Copies the init script into the VM
	command := "ssh"
	args := []string{
		"-i", filepath.Join(machine.baseDir, machine.Name, config.GetFilename(config.PrivateKeyFilename)),
		fmt.Sprintf("%s@%s", machine.Credentials.Username, machine.Network.IPAddress),
	}

	// TODO: Add stdin, stdout, stderr to os.Stdout, os.Stdin, os.Stderr
	// Run the command to create the disk
	// _, err := machine.Runner.RunCommand(command, args, osutil.WithStdin(os.Stdin), osutil.WithStdout(os.Stdout), osutil.WithStderr(os.Stderr))
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

func (machine *Machine) GetVMs() []Machine {
	return []Machine{*machine}
}

// Creates the script files from the template
func (machine *Machine) createScriptFiles() error {
	// Create the service file
	err := machine.createServiceFile()
	if err != nil {
		return err
	}

	// Create the prepare script file
	err = machine.createPreparelScriptFile()
	if err != nil {
		return err
	}

	// Create the install script file
	err = machine.createInstallScriptFile()
	if err != nil {
		return err
	}

	// Create the startup script file
	err = machine.createStartupScriptFile()
	if err != nil {
		return err
	}
	return nil
}

// Creates the service file
func (machine *Machine) createServiceFile() error {
	err := os.MkdirAll(filepath.Join(machine.baseDir, machine.Name, "bin"), 0755)
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

	err = os.WriteFile(filepath.Join(machine.baseDir, machine.Name, "bin/machina.service"), []byte(sysDSvc), 0744)
	if err != nil {
		return err
	}

	return nil
}

// Creates the install script file
func (machine *Machine) createPreparelScriptFile() error {
	// Install script
	prepareScript := `
sudo mkdir -p /etc/machina/results
sudo mv /tmp/machina/* /etc/machina
echo 'source /etc/machina/machinarc' >> $HOME/.bashrc
sudo chcon -R -t bin_t /etc/machina/machina.service
sudo cp /etc/machina/machina.service /etc/systemd/system/machina.service
sudo chmod 664 /etc/systemd/system/machina.service
sudo systemctl daemon-reload
sudo systemctl enable machina.service
sudo systemctl start machina.service
sudo chown -R %s:%s /etc/machina/*
`
	prepareScript = fmt.Sprintf(prepareScript, machine.Credentials.Username, machine.Credentials.Username)

	err := os.WriteFile(filepath.Join(machine.baseDir, machine.Name, "bin/prepare.sh"), []byte(prepareScript), 0744)
	if err != nil {
		return err
	}

	return nil
}

// Creates the install script file
func (machine *Machine) createInstallScriptFile() error {

	err := os.WriteFile(filepath.Join(machine.baseDir, machine.Name, "bin/install.sh"), []byte(machine.Scripts.Install), 0744)
	if err != nil {
		return err
	}

	return nil
}

// Create the machine startup script file
func (machine *Machine) createStartupScriptFile() error {
	// Boot script
	var mountName string

	if cfg.Hypervisor == "qemu" {
		mountName = machine.Mount.Name
	} else {
		mountName = machine.Mount.GuestPath
	}
	initScript := fmt.Sprintf(`#!/bin/bash
HOST_PATH="%s"
GUEST_PATH="%s"
if [[ "$GUEST_PATH" != "" ]]; then
	mkdir -p $GUEST_PATH
	sudo mount -t 9p $HOST_PATH $GUEST_PATH
fi

exit 0
`, mountName, machine.Mount.GuestPath)

	err := os.WriteFile(filepath.Join(machine.baseDir, machine.Name, "bin/machina"), []byte(initScript), 0744)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(machine.baseDir, machine.Name, "bin/machinarc"), []byte(machine.Scripts.Init), 0644)
	if err != nil {
		return err
	}
	return nil
}
