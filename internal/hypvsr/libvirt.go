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
	Create(machine *Machine) error
	Start(machine *Machine) error
	Stop(machine *Machine) error
	ForceStop(machine *Machine) error
	Status(machine *Machine) (string, error)
	Delete(machine *Machine) error
}

// Libvirt is a struct that represents the libvirt hypervisor
type Libvirt struct{}

// Create is a method for the libvirt hypervisor that creates an machine
func (h *Libvirt) Create(machine *Machine) error {
	// Loads the software configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	// Define the command to be executed
	command := "virt-install"

	// Validate and convert the memory
	ram, err := convertMemory(machine.Resources.Memory)
	if err != nil {
		return errors.New("invalid memory")
	}

	// Define the arguments for the execution of the command
	args := []string{
		"--connect", cfg.Connection,
		"--virt-type", "kvm",
		"--name", machine.Name,
		"--ram", ram,
		fmt.Sprintf("--vcpus=%s", machine.Resources.CPUs),
		"--os-variant", machine.Variant,
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(cfg.Directories.Instances, machine.Name, config.GetFilename(config.DiskFilename))),
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(cfg.Directories.Instances, machine.Name, config.GetFilename(config.SeedImageFilename))),
		"--import",
		"--network", fmt.Sprintf("bridge=virbr0,model=virtio,mac=%s", machine.Network.MacAddress),
		"--noautoconsole",
	}

	// Parse volume mount
	var mountCommand []string
	if machine.Mount.Name != "" {
		mountCommand = []string{
			"--filesystem",
			fmt.Sprintf("type=mount,mode=passthrough,source=%s,target=%s", machine.Mount.HostPath, machine.Mount.GuestPath),
		}
	}
	args = append(args, mountCommand...)

	// Run the command to create the machine
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return err
}

// Start is a method for the libvirt hypervisor that starts a stopped machine
func (h *Libvirt) Start(machine *Machine) error {
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
		machine.Name,
	}

	// Run the command to create the machine
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return nil
}

// Start is a method for the libvirt hypervisor that stops a running machine
func (h *Libvirt) Stop(machine *Machine) error {
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
		machine.Name,
	}

	// Run the command to create the machine
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return err
}

// ForceStop is a method for the libvirt hypervisor that force stops a running/stuck machine
func (h *Libvirt) ForceStop(machine *Machine) error {
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
		machine.Name,
	}

	// Run the command to create the machine
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return err
}

// Status is a method for the libvirt hypervisor that gets the status of an machine
func (h *Libvirt) Status(machine *Machine) (string, error) {
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
		machine.Name,
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

// Delete is a method for the libvirt hypervisor that deletes a created machine
func (h *Libvirt) Delete(machine *Machine) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Delete the machine directory
	err = os.RemoveAll(filepath.Join(cfg.Directories.Instances, machine.Name))
	if err != nil {
		return err
	}

	// Define the command to be executed
	command := "virsh"

	// Define the arguments for the execution of the command to destroy the machine
	args := []string{
		"--connect", cfg.Connection,
		"destroy",
		machine.Name,
	}

	// Run the command to create the machine
	_, err = machine.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to undefine the machine
	args = []string{
		"--connect", cfg.Connection,
		"undefine",
		machine.Name,
	}

	// Execute the command
	cmd := exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to destroy the pool of the machine
	args = []string{
		"--connect", cfg.Connection,
		"pool-destroy",
		machine.Name,
	}

	// Execute the command
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to undefine the pool of the machine
	args = []string{
		"--connect", cfg.Connection,
		"pool-undefine",
		machine.Name,
	}

	// Execute the command
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the command to be executed
	command = "ssh-keygen"

	// Define the arguments for the execution of the command to undefine the pool of the machine
	args = []string{
		"-R",
		machine.Network.IPAddress,
	}

	// Define the arguments for the execution of the command to remove the key from the SSH allowed machines
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil

}
