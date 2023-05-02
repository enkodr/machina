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

const (
	bitSize = 2048
	port    = "22"
)

type SSHClient struct {
	conn   *ssh.Client
	config *ssh.ClientConfig
}

func NewClient(host string, user string, privKeyFile string) (*SSHClient, error) {
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publicKeyFile(privKeyFile),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), cfg)
	if err != nil {
		return nil, err
	}

	return &SSHClient{
		conn:   conn,
		config: cfg,
	}, nil
}

func (c *SSHClient) RunAsSudo(command string) error {
	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(fmt.Sprintf("sudo %s", command))
}

func (c *SSHClient) RunAsUser(command string) error {
	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(fmt.Sprintf(command))
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// Generate public and private rsa keys
func GenerateNewSSHKeys() ([]byte, []byte, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

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
	publicKeyBytes := ssh.MarshalAuthorizedKey(pub)

	return privateKeyPEM, publicKeyBytes, nil
}

// Check if host is responding
func IsResponding(ip string) bool {
	timeout := time.Second * 5
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return conn != nil
}
