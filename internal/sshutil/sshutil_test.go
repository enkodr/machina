package sshutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
)

func sessionHandler(s ssh.Session) {
	io.WriteString(s, "Hello world\n")
}

const serverAddr = "127.0.0.1"

func startServer(port string) {
	s := &ssh.Server{
		Addr:    fmt.Sprintf("%s:%s", serverAddr, port),
		Handler: sessionHandler,
	}
	defer s.Close()

	log.Fatal(s.ListenAndServe())
}

// TestNewClient tests the NewClient function
func TestNewClient(t *testing.T) {
	// Test creating a new SSHClient
	port = "2222"
	go startServer(port)

	tempDir := t.TempDir()
	privKeyFile := filepath.Join(tempDir, "test_key.pem")

	// Create a temporary private key file for testing
	privKeyData := generatePrivateKeyPEM(t)
	err := ioutil.WriteFile(privKeyFile, privKeyData, 0600)
	assert.Nil(t, err)
	defer os.Remove(privKeyFile)

	// Test creating a new SSHClient
	client, err := NewClient(serverAddr, "", privKeyFile)
	assert.NotNil(t, client)
	assert.Nil(t, err)
}

// TestRunAsSudo tests the RunAsSudo function
func TestRunAsSudo(t *testing.T) {
	// Test creating a new SSHClient
	port = "2223"
	go startServer(port)

	tempDir := t.TempDir()
	privKeyFile := filepath.Join(tempDir, "test_key.pem")

	// Create a temporary private key file for testing
	privKeyData := generatePrivateKeyPEM(t)
	err := ioutil.WriteFile(privKeyFile, privKeyData, 0600)
	assert.Nil(t, err)
	defer os.Remove(privKeyFile)

	// Create a new SSHClient
	client, err := NewClient(serverAddr, "", privKeyFile)
	assert.NotNil(t, client)
	assert.Nil(t, err)

	// Test running a command as sudo
	err = client.RunAsSudo("ls")
	assert.Nil(t, err)
}

// TestRunAsUser tests the RunAsUser function
func TestRunAsUser(t *testing.T) {
	// Test creating a new SSHClient
	port = "2224"
	go startServer(port)

	tempDir := t.TempDir()
	privKeyFile := filepath.Join(tempDir, "test_key.pem")

	// Create a temporary private key file for testing
	privKeyData := generatePrivateKeyPEM(t)
	err := ioutil.WriteFile(privKeyFile, privKeyData, 0600)
	assert.Nil(t, err)
	defer os.Remove(privKeyFile)

	// Create a new SSHClient
	client, err := NewClient(serverAddr, "", privKeyFile)
	assert.NotNil(t, client)
	assert.Nil(t, err)

	// Test running a command as sudo
	err = client.RunAsUser("ls")
	assert.Nil(t, err)
}

// TestGenerateNewSSHKeys tests the GenerateNewSSHKeys function
func TestGenerateNewSSHKeys(t *testing.T) {
	// Test generating new SSH keys
	privateKeyPEM, publicKeyBytes, err := GenerateNewSSHKeys()
	assert.NotNil(t, privateKeyPEM)
	assert.NotNil(t, publicKeyBytes)
	assert.Nil(t, err)
}

// Helper function to generate a temporary private key file
func generatePrivateKeyPEM(t *testing.T) []byte {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	assert.Nil(t, err)

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   privateKeyBytes,
		},
	)

	return privateKeyPEM
}
