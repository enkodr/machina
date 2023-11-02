package sshutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
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
