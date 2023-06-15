package osutil

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"os"
	"strings"
)

// Checksum checks if the checksum of a file matches the specified checksum
func Checksum(path, checksum string) bool {
	var hasher hash.Hash
	// Split the checksum into the algorithm and the hash
	sha := strings.Split(checksum, ":")
	// Return false if the algorithm is not specified
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
