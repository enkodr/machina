package vm

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/enkodr/machina/internal/config"
)

type Libvirt struct{}

func (h *Libvirt) Create(vm *VMConfig) error {
	cfg := config.LoadConfig()
	command := "virt-install"
	ram, err := convertMemory(vm.Specs.Memory)
	if err != nil {
		return errors.New("invalid memory")
	}
	args := []string{
		"--connect", cfg.Connection,
		"--virt-type", "kvm",
		"--name", vm.Name,
		"--ram", ram,
		fmt.Sprintf("--vcpus=%s", vm.Specs.CPUs),
		"--os-variant", vm.Variant,
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(cfg.Directories.Instances, vm.Name, "disk.img")),
		"--disk", fmt.Sprintf("path=%s,device=disk", filepath.Join(cfg.Directories.Instances, vm.Name, "seed.img")),
		"--import",
		"--network", fmt.Sprintf("bridge=virbr0,model=virtio,mac=%s", vm.Network.MacAddress),
		"--noautoconsole",
	}

	args = append(args, parseLibvirtMounts(vm.Mount)...)

	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err = cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func parseLibvirtMounts(mount Mount) []string {
	if mount.Name == "" {
		return []string{}
	}
	mountCommand := []string{
		"--filesystem",
		fmt.Sprintf("type=mount,mode=passthrough,source=%s,target=%s", mount.HostPath, mount.GuestPath),
		// fmt.Sprintf("%s,%s", path, m.Path),
	}

	return mountCommand
}

func (h *Libvirt) Start(vm *VMConfig) error {
	cfg := config.LoadConfig()
	command := "virsh"
	args := []string{
		"--connect", cfg.Connection,
		"start",
		vm.Name,
	}
	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (h *Libvirt) Stop(vm *VMConfig) error {
	cfg := config.LoadConfig()
	command := "virsh"
	args := []string{
		"--connect", cfg.Connection,
		"shutdown",
		vm.Name,
	}

	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err := cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (h *Libvirt) Status(vm *VMConfig) (string, error) {
	cfg := config.LoadConfig()
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

func (h *Libvirt) Delete(vm *VMConfig) error {
	cfg := config.LoadConfig()
	command := "virsh"
	args := []string{
		"--connect", cfg.Connection,
		"destroy",
		vm.Name,
	}

	logFile, _ := os.OpenFile(filepath.Join(cfg.Directories.Instances, vm.Name, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()
	cmd := exec.Command(command, args...)
	cmd.Stdout = logFile

	err := cmd.Start()
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

	return os.RemoveAll(filepath.Join(cfg.Directories.Instances, vm.Name))

}
