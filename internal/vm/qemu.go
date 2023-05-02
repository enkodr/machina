package vm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/enkodr/machina/internal/config"
)

type Qemu struct{}

func (h *Qemu) Create(vm *VMConfig) error {
	return h.Start(vm)
}

func (h *Qemu) Start(vm *VMConfig) error {
	command := "qemu-system-x86_64"
	cfg := config.LoadConfig()
	dir := filepath.Join(cfg.Directories.Instances, vm.Name)
	args := []string{
		"-machine", fmt.Sprintf("accel=%s,type=q35", getHypervisor()),
		"-cpu", "host",
		"-smp", vm.Specs.CPUs,
		"-m", vm.Specs.Memory,
		"-nographic",
		"-netdev", fmt.Sprintf("bridge,id=%s,br=virbr0", vm.Network.NicName),
		"-device", fmt.Sprintf("virtio-net-pci,netdev=%s,id=virtnet0,mac=%s", vm.Network.NicName, vm.Network.MacAddress),
		"-pidfile", fmt.Sprintf("%s/vm.pid", dir),
		"-drive", fmt.Sprintf("if=virtio,format=qcow2,file=%s/disk.img", dir),
		"-drive", fmt.Sprintf("if=virtio,format=raw,file=%s/seed.img", dir),
	}

	args = append(args, parseQemuMounts(vm.Mount)...)
	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (h *Qemu) Stop(vm *VMConfig) error {
	cfg := config.LoadConfig()
	command := "ssh"
	args := []string{
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(cfg.Directories.Instances, vm.Name, "id_rsa"),
		fmt.Sprintf("%s@%s", vm.Credentials.Username, vm.Network.IPAddress),
		"sudo", "shutdown", "now",
	}
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (h *Qemu) Status(vm *VMConfig) (string, error) {
	cfg := config.LoadConfig()
	if _, err := os.Stat(filepath.Join(cfg.Directories.Instances, vm.Name, "vm.pid")); os.IsNotExist(err) {
		return "shut off", nil
	}
	return "running", nil
}

func (h *Qemu) Delete(vm *VMConfig) error {
	cfg := config.LoadConfig()

	status, _ := h.Status(vm)
	if status == "running" {
		h.Stop(vm)
	}

	command := "ssh-keygen"
	args := []string{
		"-R",
		vm.Network.IPAddress,
	}
	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join(cfg.Directories.Instances, vm.Name))
}

func parseQemuMounts(mount Mount) []string {
	if mount.Name == "" {
		return []string{}
	}
	home, _ := os.UserHomeDir()
	path := strings.Replace(mount.HostPath, "~", home, -1)
	mountCommand := []string{
		"-fsdev",
		fmt.Sprintf("local,security_model=passthrough,id=fsdev%d,path=%s", 0, path),
		"--device",
		fmt.Sprintf("virtio-9p-pci,id=fs%d,fsdev=fsdev%d,mount_tag=%s", 0, 0, mount.Name),
	}
	return mountCommand
}

// Get Hypervisor driver
func getHypervisor() string {
	driver := "kvm"

	if runtime.GOOS == "darwin" {
		driver = "hvf"
	}

	return driver
}
