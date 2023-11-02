package sshutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	bitSize = 2048
	port    = "22"
)

// SSHClient is a wrapper around the ssh.Client
type SSHClient struct {
	conn   *ssh.Client
	config *ssh.ClientConfig
}

// NewClient creates a new SSHClient
func NewClient(host string, user string, privKeyFile string) (*SSHClient, error) {
	// Create the ssh config
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publicKeyFile(privKeyFile),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// Start the ssh connection
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), cfg)
	if err != nil {
		return nil, err
	}

	// Return the ssh client
	return &SSHClient{
		conn:   conn,
		config: cfg,
	}, nil
}

// RunAsSudo runs a command as sudo on the remote host
func (c *SSHClient) RunAsSudo(command string) error {
	// Create a new session
	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Run the command
	return session.Run(fmt.Sprintf("sudo %s", command))
}

// RunAsUser runs a command as the user on the remote host
func (c *SSHClient) RunAsUser(command string) error {
	// Create a new session
	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Run the command
	return session.Run(fmt.Sprintf(command))
}

// publicKeyFile reads the private key file and returns the ssh.AuthMethod
func publicKeyFile(file string) ssh.AuthMethod {
	// Read the private key file
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	// Parse the private key
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	// Return the ssh.AuthMethod
	return ssh.PublicKeys(key)
}

// GenerateNewSSHKeys generates a new SSH key pair
func GenerateNewSSHKeys() ([]byte, []byte, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

	// Generate private key PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   privateKeyBytes,
		},
	)

	// Generate public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	// Generate public key bytes
	publicKeyBytes := ssh.MarshalAuthorizedKey(pub)

	// Return the private key PEM and public key bytes
	return privateKeyPEM, publicKeyBytes, nil
}

// Check if host is responding
func IsResponding(ip string) bool {
	// Set timeout to 5 seconds
	timeout := time.Second * 5

	// Try to connect to the host
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	// Return true if the connection is established
	return conn != nil
}
