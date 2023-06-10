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

type Hypervisor interface {
	Create(vm *Machine) error
	Start(vm *Machine) error
	Stop(vm *Machine) error
	ForceStop(vm *Machine) error
	Status(vm *Machine) (string, error)
	Delete(vm *Machine) error
}

type Libvirt struct{}

func (h *Libvirt) Create(vm *Machine) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	command := "virt-install"
	ram, err := convertMemory(vm.Resources.Memory)
	if err != nil {
		return errors.New("invalid memory")
	}
	args := []string{
		"--connect", cfg.Connection,
		"--virt-type", "kvm",
		"--name", vm.Name,
		"--ram", ram,
		fmt.Sprintf("--vcpus=%s", vm.Resources.CPUs),
		"--os-variant", vm.Variant,
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(cfg.Directories.Machines, vm.Name, config.GetFilename(config.DiskFilename))),
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(cfg.Directories.Machines, vm.Name, config.GetFilename(config.SeedImageFilename))),
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

	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Machines, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err = cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (h *Libvirt) Start(vm *Machine) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	command := "virsh"
	args := []string{
		"--connect", cfg.Connection,
		"start",
		vm.Name,
	}
	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Machines, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (h *Libvirt) Stop(vm *Machine) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	command := "virsh"
	args := []string{
		"--connect", cfg.Connection,
		"shutdown",
		vm.Name,
	}

	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Machines, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err = cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (h *Libvirt) ForceStop(vm *Machine) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	command := "virsh"
	args := []string{
		"--connect", cfg.Connection,
		"destroy",
		vm.Name,
	}

	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Machines, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err = cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (h *Libvirt) Status(vm *Machine) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	command := "virsh"
	args := []string{
		"--connect", cfg.Connection,
		"domstate",
		vm.Name,
	}

	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	strOutput := strings.TrimSpace(string(output))
	return strOutput, nil
}

func (h *Libvirt) Delete(vm *Machine) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	command := "virsh"
	args := []string{
		"--connect", cfg.Connection,
		"destroy",
		vm.Name,
	}

	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Machines, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err = cmd.Start()
	if err != nil {
		return err
	}

	args = []string{
		"--connect", cfg.Connection,
		"undefine",
		vm.Name,
	}
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	args = []string{
		"--connect", cfg.Connection,
		"pool-destroy",
		vm.Name,
	}
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	args = []string{
		"--connect", cfg.Connection,
		"pool-undefine",
		vm.Name,
	}
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	command = "ssh-keygen"
	args = []string{
		"-R",
		vm.Network.IPAddress,
	}
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join(cfg.Directories.Machines, vm.Name))

}
