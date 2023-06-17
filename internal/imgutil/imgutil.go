package imgutil

import (
	"errors"
	"os"
	"strings"

	"github.com/enkodr/machina/internal/config"
)

func EnsureDirectories(cfg *config.Config) {
	// Create directories
	os.MkdirAll(cfg.Directories.Images, 0755)
	os.MkdirAll(cfg.Directories.Instances, 0755)

}

// GetFileName extracts and returns the name of the file from the URL
func GetFilenameFromURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("empty URL")
	}

	parts := strings.Split(url, "/")
	size := len(parts)

	if size < 2 {
		return "", errors.New("invalid URL format")
	}

	if size < 4 {
		return "", errors.New("no file in the URL")
	}

	filename := parts[size-1]
	if filename == "" {
		return "", errors.New("no file in the URL")
	}

	return filename, nil
}
