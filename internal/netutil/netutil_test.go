package netutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNetwork(t *testing.T) {
	// Test case 1: Not nill
	net := NewNetwork()
	assert.NotNil(t, net)

	// Test case 2: Valid Data
	// Check that the interface has the expected name
	assert.Equal(t, "virtnet", net.Ethernets.VirtNet.Name)

	// Check that the interface has the expected IP address and subnet mask
	got := net.Ethernets.VirtNet.Addresses[0]
	want := ipRange
	assert.True(t, strings.Contains(got, want))

	// Check that the interface has the expected gateway address
	baseIp := strings.Split(net.Ethernets.VirtNet.Addresses[0], "/")[0]
	want, _ = GetGatewayFromIP(baseIp)
	got = net.Ethernets.VirtNet.Gateway4
	assert.Equal(t, want, got)

	// Check that DHCP is disabled
	assert.False(t, net.Ethernets.VirtNet.DHCP4)

	// Check that the interface has a random MAC address
	assert.NotEmpty(t, net.Ethernets.VirtNet.Match.MacAddress)

	// Check that the interface has the expected nameserver addresses
	assert.Equal(t, nameservers, net.Ethernets.VirtNet.Nameservers.Addresses)
}

func TestGetGatewayFromIP(t *testing.T) {
	ip := "192.168.1.10"
	want := "192.168.1.1"

	// Test Case 1: Test with valid IP address
	got, _ := GetGatewayFromIP(ip)
	assert.Equal(t, want, got)

	// Test Case 2: Test with invalid IP address
	_, err := GetGatewayFromIP("192.168.122")
	assert.Error(t, err)
}

func TestValidateIPAddress(t *testing.T) {
	valid := "192.168.1.1"
	invalid := "192.168.1.256"

	if !ValidateIPAddress(valid) {
		t.Errorf("IP %s should be valid", valid)
	}

	if ValidateIPAddress(invalid) {
		t.Errorf("IP %s should be invalid", invalid)
	}
}

func TestGenerateIpAddress(t *testing.T) {
	ip := GenerateIPAddress()
	if !ValidateIPAddress(ip) {
		t.Errorf("failed to generate IP address")
	}
}

func TestRandomMacAddress(t *testing.T) {
	// Test Case 1: Unable to generate MacAddress
	_, err := RandomMacAddress()
	assert.NoError(t, err)
}

func TestDownload(t *testing.T) {
	// Create a test server with a custom handler
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a success status code and test file content
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Test file content")
	}))
	defer ts.Close()

	// Use the test server URL for the download
	url := ts.URL

	data, err := Download(url)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Test file content"), data)
}

func TestDownloadAndSave(t *testing.T) {
	// Create a temporary directory for file saving
	tempDir := t.TempDir()

	// Create a test server with a custom handler
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a success status code and test file content
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Test file content")
	}))
	defer ts.Close()

	// Use the test server URL and temporary directory for the download and save
	url := fmt.Sprintf("%s/file.txt", ts.URL)
	destination := tempDir

	err := DownloadAndSave(url, destination)
	assert.Nil(t, err)

	// Check if the file was saved correctly
	filePath := filepath.Join(destination, "file.txt")
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Test file content"), data)
}

func TestGetIPFromNetworkAddress(t *testing.T) {
	// Test case with a valid network address
	netAddr := "192.168.0.0/24"
	expectedIP := "192.168.0.0"
	ip, err := GetIPFromNetworkAddress(netAddr)
	assert.Nil(t, err)
	assert.Equal(t, expectedIP, ip)

	// Test case with an invalid network address
	netAddr = "192.168.0.0" // Missing CIDR notation
	ip, err = GetIPFromNetworkAddress(netAddr)
	assert.NotNil(t, err)
	assert.Equal(t, "", ip)

	// Test case with an empty network address
	netAddr = "" // Missing CIDR notation
	ip, err = GetIPFromNetworkAddress(netAddr)
	assert.NotNil(t, err)
	assert.Equal(t, "", ip)
}
