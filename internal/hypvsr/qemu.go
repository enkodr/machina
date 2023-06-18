package hypvsr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/enkodr/machina/internal/config"
)

type Qemu struct{}

func (h *Qemu) Create(vm *Machine) error {
	return h.Start(vm)
}

func (h *Qemu) Start(vm *Machine) error {
	command := "qemu-system-x86_64"
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	dir := filepath.Join(cfg.Directories.Instances, vm.Name)
	args := []string{
		"-machine", fmt.Sprintf("accel=%s,type=q35", getHypervisor()),
		"-cpu", "host",
		"-smp", vm.Resources.CPUs,
		"-m", vm.Resources.Memory,
		"-nographic",
		"-netdev", fmt.Sprintf("bridge,id=%s,br=virbr0", vm.Network.NicName),
		"-device", fmt.Sprintf("virtio-net-pci,netdev=%s,id=virtnet0,mac=%s", vm.Network.NicName, vm.Network.MacAddress),
		"-pidfile", fmt.Sprintf("%s/vm.pid", dir),
		"-drive", fmt.Sprintf("if=virtio,format=qcow2,file=%s/disk.img", dir),
		"-drive", fmt.Sprintf("if=virtio,format=raw,file=%s/seed.img", dir),
	}

	var mountCommand []string
	if vm.Mount.Name != "" {
		mountCommand = []string{
			"-fsdev",
			fmt.Sprintf("local,security_model=passthrough,id=fsdev%d,path=%s", 0, vm.Mount.HostPath),
			"--device",
			fmt.Sprintf("virtio-9p-pci,id=fs%d,fsdev=fsdev%d,mount_tag=%s", 0, 0, vm.Mount.Name),
		}
	}
	args = append(args, mountCommand...)
	cmd := exec.Command(command, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (h *Qemu) Stop(vm *Machine) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	command := "ssh"
	args := []string{
		"-o StrictHostKeyChecking=no",
		"-i", filepath.Join(cfg.Directories.Instances, vm.Name, config.GetFilename(config.PrivateKeyFilename)),
		fmt.Sprintf("%s@%s", vm.Credentials.Username, vm.Network.IPAddress),
		"sudo", "shutdown", "now",
	}
	cmd := exec.Command(command, args...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (h *Qemu) ForceStop(vm *Machine) error {
	command := "kill"
	args := []string{
		"-9",
		h.getPID(vm),
	}
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (h *Qemu) Status(vm *Machine) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(filepath.Join(cfg.Directories.Instances, vm.Name, config.GetFilename(config.PIDFilename))); os.IsNotExist(err) {
		return "shut off", nil
	}
	return "running", nil
}

func (h *Qemu) Delete(vm *Machine) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

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
	err = cmd.Start()
	if err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join(cfg.Directories.Instances, vm.Name))
}

// Get Hypervisor driver
func getHypervisorDriver() string {
	driver := "kvm"

	if runtime.GOOS == "darwin" {
		driver = "hvf"
	}

	return driver
}

func (h *Qemu) getPID(vm *Machine) string {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "err"
	}
	data, _ := os.ReadFile(filepath.Join(cfg.Directories.Instances, vm.Name, config.GetFilename(config.PIDFilename)))
	return string(data)
}
