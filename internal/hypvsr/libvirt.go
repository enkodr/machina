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
	Create(instance *Instance) error
	Start(instance *Instance) error
	Stop(instance *Instance) error
	ForceStop(instance *Instance) error
	Status(instance *Instance) (string, error)
	Delete(instance *Instance) error
}

// Libvirt is a struct that represents the libvirt hypervisor
type Libvirt struct{}

// Create is a method for the libvirt hypervisor that creates an instance
func (h *Libvirt) Create(instance *Instance) error {
	// Loads the software configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	// Define the command to be executed
	command := "virt-install"

	// Validate and convert the memory
	ram, err := convertMemory(instance.Resources.Memory)
	if err != nil {
		return errors.New("invalid memory")
	}

	// Define the arguments for the execution of the command
	args := []string{
		"--connect", cfg.Connection,
		"--virt-type", "kvm",
		"--name", instance.Name,
		"--ram", ram,
		fmt.Sprintf("--vcpus=%s", instance.Resources.CPUs),
		"--os-variant", instance.Variant,
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.DiskFilename))),
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(instance.baseDir, instance.Name, config.GetFilename(config.SeedImageFilename))),
		"--import",
		"--network", fmt.Sprintf("bridge=virbr0,model=virtio,mac=%s", instance.Network.MacAddress),
		"--noautoconsole",
	}

	// Parse volume mount
	var mountCommand []string
	if instance.Mount.Name != "" {
		mountCommand = []string{
			"--filesystem",
			fmt.Sprintf("type=mount,mode=passthrough,source=%s,target=%s", instance.Mount.HostPath, instance.Mount.GuestPath),
		}
	}
	args = append(args, mountCommand...)

	// Run the command to create the instance
	_, err = instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return err
}

// Start is a method for the libvirt hypervisor that starts a stopped instance
func (h *Libvirt) Start(instance *Instance) error {
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
		instance.Name,
	}

	// Run the command to create the instance
	_, err = instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return nil
}

// Start is a method for the libvirt hypervisor that stops a running instance
func (h *Libvirt) Stop(instance *Instance) error {
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
		instance.Name,
	}

	// Run the command to create the instance
	_, err = instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return err
}

// ForceStop is a method for the libvirt hypervisor that force stops a running/stuck instance
func (h *Libvirt) ForceStop(instance *Instance) error {
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
		instance.Name,
	}

	// Run the command to create the instance
	_, err = instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	return err
}

// Status is a method for the libvirt hypervisor that gets the status of an instance
func (h *Libvirt) Status(instance *Instance) (string, error) {
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
		instance.Name,
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
func (h *Libvirt) Delete(instance *Instance) error {
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
		instance.Name,
	}

	// Run the command to create the instance
	_, err = instance.Runner.RunCommand(command, args)
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to undefine the instance
	args = []string{
		"--connect", cfg.Connection,
		"undefine",
		instance.Name,
	}

	// Execute the command
	cmd := exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to destroy the pool of the instance
	args = []string{
		"--connect", cfg.Connection,
		"pool-destroy",
		instance.Name,
	}

	// Execute the command
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the arguments for the execution of the command to undefine the pool of the instance
	args = []string{
		"--connect", cfg.Connection,
		"pool-undefine",
		instance.Name,
	}

	// Execute the command
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Define the command to be executed
	command = "ssh-keygen"

	// Define the arguments for the execution of the command to undefine the pool of the instance
	args = []string{
		"-R",
		instance.Network.IPAddress,
	}

	// Define the arguments for the execution of the command to remove the key from the SSH allowed machines
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Delete the instance directory
	return os.RemoveAll(filepath.Join(instance.baseDir, instance.Name))

}
