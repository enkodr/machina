package osutil

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Runner interface {
	RunCommand(command string, args []string, options ...Option) (string, error)
}

type Option func(*commandOptions)

type commandOptions struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// WithStdin sets the stdin for the command.
func WithStdin(stdin io.Reader) Option {
	return func(opts *commandOptions) {
		opts.stdin = stdin
	}
}

// WithStdout sets the stdout for the command.
func WithStdout(stdout io.Writer) Option {
	return func(opts *commandOptions) {
		opts.stdout = stdout
	}
}

// WithStderr sets the stderr for the command.
func WithStderr(stderr io.Writer) Option {
	return func(opts *commandOptions) {
		opts.stderr = stderr
	}
}

// CommandRunner is the implementation of the Runner interface that runs OS commands.
type CommandRunner struct{}

// RunCommand runs the given OS command with the specified options, and returns the output as a string,
// along with any error that occurred.
func (c CommandRunner) RunCommand(command string, args []string, options ...Option) (string, error) {
	opts := &commandOptions{}
	for _, opt := range options {
		opt(opts)
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = opts.stdin
	cmd.Stdout = opts.stdout
	cmd.Stderr = opts.stderr

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

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
