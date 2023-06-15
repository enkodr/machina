package hypvsr

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/enkodr/machina/internal/config"
)

// Hypervisor is an interface for interacting with the hypervisor
type Hypervisor interface {
	Create(vm *Instance) error
	Start(vm *Instance) error
	Stop(vm *Instance) error
	ForceStop(vm *Instance) error
	Status(vm *Instance) (string, error)
	Delete(vm *Instance) error
}

// Libvirt is a struct that represents the libvirt hypervisor
type Libvirt struct{}

// Create is a method for the libvirt hypervisor that creates an instance
func (h *Libvirt) Create(vm *Instance) error {
	// Loads the software configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	// Define the command to be executed
	command := "virt-install"

	// Validate and convert the memory
	ram, err := convertMemory(vm.Resources.Memory)
	if err != nil {
		return errors.New("invalid memory")
	}

	// Define the arguments for the execution of the command
	args := []string{
		"--connect", cfg.Connection,
		"--virt-type", "kvm",
		"--name", vm.Name,
		"--ram", ram,
		fmt.Sprintf("--vcpus=%s", vm.Resources.CPUs),
		"--os-variant", vm.Variant,
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(cfg.Directories.Instances, vm.Name, config.GetFilename(config.DiskFilename))),
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(cfg.Directories.Instances, vm.Name, config.GetFilename(config.SeedImageFilename))),
		"--import",
		"--network", fmt.Sprintf("bridge=virbr0,model=virtio,mac=%s", vm.Network.MacAddress),
		"--noautoconsole",
	}

	// Parse volume mount
	var mountCommand []string
	if vm.Mount.Name != "" {
		mountCommand = []string{
			"--filesystem",
			fmt.Sprintf("type=mount,mode=passthrough,source=%s,target=%s", vm.Mount.HostPath, vm.Mount.GuestPath),
		}
	}
	args = append(args, mountCommand...)

	// Logs the output of the command into the instance logfile
	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()

	// Execute the command
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile
	err = cmd.Start()
	if err != nil {
		return err
	}
	return err
}

// Start is a method for the libvirt hypervisor that starts a stopped instance
func (h *Libvirt) Start(vm *Instance) error {
	// Loads the software configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Define the command to be executed
	command := "virsh"

	// Define the arguments for the execution of the command
	args := []string{
		"--connect", cfg.Connection,
		"start",
		vm.Name,
	}

	// Logs the output of the command into the instance logfile
	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()

	// Execute the command
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

// Start is a method for the libvirt hypervisor that stops a running instance
func (h *Libvirt) Stop(vm *Instance) error {
	// Loads the software configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Define the command to be executed
	command := "virsh"

	// Define the arguments for the execution of the command
	args := []string{
		"--connect", cfg.Connection,
		"shutdown",
		vm.Name,
	}

	// Logs the output of the command into the instance logfile
	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()

	// Execute the command
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile
	err = cmd.Start()
	if err != nil {
		return err
	}
	return err
}

// ForceStop is a method for the libvirt hypervisor that force stops a running/stuck instance
func (h *Libvirt) ForceStop(vm *Instance) error {
	// Loads the software configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Define the command to be executed
	command := "virsh"

	// Define the arguments for the execution of the command
	args := []string{
		"--connect", cfg.Connection,
		"destroy",
		vm.Name,
	}

	// Logs the output of the command into the instance logfile
	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()

	// Execute the command
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile
	err = cmd.Start()
	if err != nil {
		return err
	}
	return err
}

// Status is a method for the libvirt hypervisor that gets the status of an instance
func (h *Libvirt) Status(vm *Instance) (string, error) {
	// Loads the software configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	// Define the command to be executed
	command := "virsh"

	// Define the arguments for the execution of the command
	args := []string{
		"--connect", cfg.Connection,
		"domstate",
		vm.Name,
	}

	// Execute the command
	cmd := exec.Command(command, args...)

	// Runs the command and gets the combined outpu
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	strOutput := strings.TrimSpace(string(output))
	return strOutput, nil
}

// Delete is a method for the libvirt hypervisor that deletes a created instance
func (h *Libvirt) Delete(vm *Instance) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Define the command to be executed
	command := "virsh"

	// Define the arguments for the execution of the command to destroy the instance
	args := []string{
		"--connect", cfg.Connection,
		"destroy",
		vm.Name,
	}

	// Logs the output of the command into the instance logfile
	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()

	// Execute the command
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to undefine the instance
	args = []string{
		"--connect", cfg.Connection,
		"undefine",
		vm.Name,
	}

	// Execute the command
	cmd = exec.Command(command, args...)
	cmd.Stdout = logFile
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to destroy the pool of the instance
	args = []string{
		"--connect", cfg.Connection,
		"pool-destroy",
		vm.Name,
	}

	// Execute the command
	cmd = exec.Command(command, args...)
	cmd.Stdout = logFile
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to undefine the pool of the instance
	args = []string{
		"--connect", cfg.Connection,
		"pool-undefine",
		vm.Name,
	}

	// Execute the command
	cmd = exec.Command(command, args...)
	cmd.Stdout = logFile
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the command to be executed
	command = "ssh-keygen"

	// Define the arguments for the execution of the command to undefine the pool of the instance
	args = []string{
		"-R",
		vm.Network.IPAddress,
	}

	// Define the arguments for the execution of the command to remove the key from the SSH allowed machines
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Delete the instance directory
	return os.RemoveAll(filepath.Join(cfg.Directories.Instances, vm.Name))

}
