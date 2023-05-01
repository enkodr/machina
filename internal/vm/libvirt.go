package vm

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

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
		// "--filesystem", fmt.Sprintf("type=mount,mode=passthrough,source=%s,target=%s", m.Mounts[0].HostPath, m.Mounts[0].Path),
		"--network", fmt.Sprintf("bridge=virbr0,model=virtio,mac=%s", vm.Network.MacAddress),
		"--noautoconsole",
	}

	cmd := exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}
	return err
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

func (h *Libvirt) Start(vm *VMConfig) error {
	cfg := config.LoadConfig()
	command := "virsh"
	args := []string{
		"--connect", fmt.Sprintf("%s", cfg.Connection),
		"start",
		vm.Name,
	}

	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (h *Libvirt) Stop(vm *VMConfig) error {
	cfg := config.LoadConfig()
	command := "virsh"
	args := []string{
		"--connect", fmt.Sprintf("%s", cfg.Connection),
		"shutdown",
		vm.Name,
	}

	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (h *Libvirt) Status(vm *VMConfig) error {
	command := "virt-install"
	args := []string{}

	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (h *Libvirt) Delete(vm *VMConfig) error {
	cfg := config.LoadConfig()
	command := "virsh"
	args := []string{
		"--connect", fmt.Sprintf("%s", cfg.Connection),
		"destroy",
		vm.Name,
	}
	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	args = []string{
		"--connect", fmt.Sprintf("%s", cfg.Connection),
		"undefine",
		vm.Name,
	}
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	args = []string{
		"--connect", fmt.Sprintf("%s", cfg.Connection),
		"pool-destroy",
		vm.Name,
	}
	cmd = exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	args = []string{
		"--connect", fmt.Sprintf("%s", cfg.Connection),
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

	return os.RemoveAll(filepath.Join(cfg.Directories.Instances, vm.Name))

}
