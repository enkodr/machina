package sshutil

import (
	"fmt"
	"io"
	"log"
	"testing"

	sshsrv "github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
)

func sessionHandler(s sshsrv.Session) {
	io.WriteString(s, "Hello world\n")
}

const serverAddr = "127.0.0.1"

func startServer(port string) {
	s := &sshsrv.Server{
		Addr:    fmt.Sprintf("%s:%s", serverAddr, port),
		Handler: sessionHandler,
	}
	defer s.Close()

	log.Fatal(s.ListenAndServe())
}

// TestGenerateNewSSHKeys tests the GenerateNewSSHKeys function
func TestGenerateNewSSHKeys(t *testing.T) {
	// Test generating new SSH keys
	privateKeyPEM, publicKeyBytes, err := GenerateNewSSHKeys()
	assert.NotNil(t, privateKeyPEM)
	assert.NotNil(t, publicKeyBytes)
	assert.Nil(t, err)
}
