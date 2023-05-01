package imgutil

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/osutil"
)

func EnsureDirectories(name string) {
	cfg := config.LoadConfig()
	osutil.MkDir(cfg.Directories.Images)
	osutil.MkDir(filepath.Join(cfg.Directories.Instances, name))
}

// GetFileName extracts and returns the name of the file from the URL
func GetFilenameFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")
	size := len(parts)
	if size < 4 {
		return "", errors.New("no file in the url")
	}
	return parts[size-1], nil
}
