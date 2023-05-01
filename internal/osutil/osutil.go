package osutil

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"os"
	"strings"
)

// create a directory if it doesn't exist yet
func MkDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

// CheckSum validates the file checksum agains the one provided
func Checksum(path, checksum string) bool {
	// Get the hash algorithm from what is specified for the machine
	var hasher hash.Hash
	sha := strings.Split(checksum, ":")
	if len(sha) == 1 {
		return false
	}

	// Select the correct halgorithm based on the
	switch sha[0] {
	case "sha256":
		hasher = sha256.New()
	case "sha512":
		hasher = sha512.New()
	default:
		return false
	}

	// Get the hash from the file
	s, err := os.ReadFile(path)
	hasher.Write(s)
	if err != nil {
		return false
	}

	// Check and return if the match
	return sha[1] == hex.EncodeToString(hasher.Sum(nil))
}
