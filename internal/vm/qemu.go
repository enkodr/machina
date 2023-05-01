package vm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/enkodr/machina/internal/config"
)

type Qemu struct{}

func (h *Qemu) Create(vm *VMConfig) error {
	return nil
}

func (h *Qemu) Start(vm *VMConfig) error {
	command := "qemu-system-x86_64"
	cfg := config.LoadConfig()
	dir := filepath.Join(cfg.Directories.Instances, vm.Name)
	args := []string{
		"-machine", "accel=kvm,type=q35",
		"-cpu", "host",
		"-smp", vm.Specs.CPUs,
		"-m", vm.Specs.Memory,
		"-nographic",
		"-netdev", fmt.Sprintf("bridge,id=%s,br=virbr0", vm.Network.NicName),
		"-device", fmt.Sprintf("virtio-net-pci,netdev=%s,id=virtnet0,mac=%s", vm.Network.NicName, vm.Network.MacAddress),
		parseQemuMounts(vm.Mount),
		"-pidfile", fmt.Sprintf("%s/vm.pid", dir),
		"-drive", fmt.Sprintf("if=virtio,format=qcow2,file=%s/disk.img", dir),
		"-drive", fmt.Sprintf("if=virtio,format=raw,file=%s/seed.img", dir),
	}
	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (h *Qemu) Stop(vm *VMConfig) error {
	return nil
}

func (h *Qemu) Status(vm *VMConfig) (string, error) {
	return "", nil
}

func (h *Qemu) Delete(vm *VMConfig) error {
	cfg := config.LoadConfig()

	return os.RemoveAll(filepath.Join(cfg.Directories.Instances, vm.Name))
}

func parseQemuMounts(mount Mount) string {
	home, _ := os.UserHomeDir()
	path := strings.Replace(mount.HostPath, "~", home, -1)
	mountCommand := []string{
		"-fsdev",
		fmt.Sprintf("local,security_model=passthrough,id=fsdev%d,path=%s", 0, path),
		"--device",
		fmt.Sprintf("virtio-9p-pci,id=fs%d,fsdev=fsdev%d,mount_tag=%s", 0, 0, mount.Name),
	}
	return strings.Join(mountCommand, " ")
}
