package imgutil

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/osutil"
)

func EnsureDirectories() error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	osutil.MkDir(cfg.Directories.Images)
	osutil.MkDir(filepath.Join(cfg.Directories.Machines))
	osutil.MkDir(filepath.Join(cfg.Directories.Clusters))

	return nil
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
