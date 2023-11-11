package path

import (
	"os"
	"path/filepath"
)

type Filename int

const (
	DatabaseFile Filename = iota
	NetworkFile
	UserdataFile
	PrivateKeyFile
	InstanceFile
	SeedImageFile
	DiskFile
	PIDFile
	ImagesDir
	InstancesDir
	ResultsDir
)

var (
	baseDir = ".local/share/machina"
)

func GetPath(fn Filename) string {
	home, _ := os.UserHomeDir()
	switch fn {
	case DatabaseFile:
		return filepath.Join(home, baseDir, "machina.db")
	case NetworkFile:
		return filepath.Join(home, baseDir, "instances", "%s", "network.cfg")
	case UserdataFile:
		return filepath.Join(home, baseDir, "instances", "%s", "userdata.yaml")
	case PrivateKeyFile:
		return filepath.Join(home, baseDir, "instances", "%s", "id_rsa")
	case InstanceFile:
		return filepath.Join(home, baseDir, "instances", "%s", "instance.yaml")
	case SeedImageFile:
		return filepath.Join(home, baseDir, "instances", "%s", "seed.img")
	case DiskFile:
		return filepath.Join(home, baseDir, "instances", "%s", "disk.img")
	case PIDFile:
		return filepath.Join(home, baseDir, "instances", "%s", "vm.pid")
	case ImagesDir:
		return filepath.Join(home, baseDir, "images")
	case InstancesDir:
		return filepath.Join(home, baseDir, "instances")
	case ResultsDir:
		return filepath.Join(home, baseDir, "results")
	}

	return ""
}
