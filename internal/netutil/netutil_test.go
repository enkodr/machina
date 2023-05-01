package netutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/file.txt")
	}))
	defer testServer.Close()

	// Test case 1: Valid download
	url := fmt.Sprintf("%s/file.txt", testServer.URL)
	dwnlData, err := Download(url)
	assert.NoError(t, err)

	// Test case 2: Validate downloaded data
	fileData, err := os.ReadFile("testdata/file.txt")
	assert.NoError(t, err)
	assert.Equal(t, fileData, dwnlData)
}

func TestDownloadNotFound(t *testing.T) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer testServer.Close()

	// Test case 1: Valid download
	url := fmt.Sprintf("%s/file.txt", testServer.URL)
	_, err := Download(url)
	assert.Error(t, err)

}

func TestDownloadAndSave(t *testing.T) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/file.txt")
	}))
	defer testServer.Close()
	url := fmt.Sprintf("%s/file.txt", testServer.URL)

	destination := "/tmp"
	err := DownloadAndSave(url, destination)
	assert.NoError(t, err)

	// Verify that the file was downloadADownloadAndSaveed to the correct location
	filePath := filepath.Join(destination, "file.txt")
	_, err = os.Stat(filePath)
	assert.False(t, os.IsNotExist(err))

	// Clean up the test file
	os.Remove(filePath)

	// Test case 2: Invalid URL
	err = DownloadAndSave("invalidurl", destination)
	assert.Error(t, err)

	// Test case 3: Valid URL but with no file specified
	err = DownloadAndSave(testServer.URL, destination)
	assert.Error(t, err)

	// Test case 4: Invalid destination
	destination = "invalidpath"
	err = DownloadAndSave(url, destination)
	assert.Error(t, err)

}

func TestGetIPFromNetworkAddress(t *testing.T) {
	ips := []struct {
		input string
		want  string
	}{
		{"192.168.0.1/24", "192.168.0.1"},
		{"10.0.0.1/16", "10.0.0.1"},
		{"172.16.0.1/12", "172.16.0.1"},
		{"", ""},
		{"192.168.0.1", "192.168.0.1"}, // Test without network prefix
	}

	for _, ip := range ips {
		got := GetIPFromNetworkAddress(ip.input)
		assert.Equal(t, got, ip.want)
	}
}
