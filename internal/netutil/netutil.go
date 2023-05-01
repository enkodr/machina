package netutil

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	rnd "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/enkodr/machina/internal/imgutil"
)

type Network struct {
	Ethernets Ethernets `yaml:"ethernets"`
	Version   int       `yaml:"version"`
}

type Ethernets struct {
	VirtNet VirtNet `yaml:"virtnet"`
}

type VirtNet struct {
	Name        string      `yaml:"set-name"`
	Addresses   []string    `yaml:"addresses"`
	DHCP4       bool        `yaml:"dhcp4"`
	Gateway4    string      `yaml:"gateway4"`
	Match       Match       `yaml:"match"`
	Nameservers Nameservers `yaml:"nameservers"`
}

type Match struct {
	MacAddress string `yaml:"macaddress"`
}

type Nameservers struct {
	Addresses []string `yaml:"addresses"`
}

var (
	ipRange     = "192.168.122"
	nameservers = []string{"1.1.1.1", "8.8.8.8"}
)

func NewNetwork() *Network {
	net := &Network{}
	ipAddress := GenerateIPAddress()
	gwAddress, _ := GetGatewayFromIP(ipAddress)
	macAddress, _ := RandomMacAddress()
	net.Version = 2
	net.Ethernets.VirtNet.Name = "virtnet"
	net.Ethernets.VirtNet.Addresses = append(net.Ethernets.VirtNet.Addresses, fmt.Sprintf("%s/24", ipAddress))
	net.Ethernets.VirtNet.Gateway4 = gwAddress
	net.Ethernets.VirtNet.DHCP4 = false
	net.Ethernets.VirtNet.Match.MacAddress = macAddress
	net.Ethernets.VirtNet.Nameservers.Addresses = append(net.Ethernets.VirtNet.Nameservers.Addresses, nameservers...)
	return net
}

// getGatewayFromIP will return the gateway from the IPv4 addres
//
//	Example:
//		IP: 192.168.122.100
//		GW? 192.168.122.1
func GetGatewayFromIP(ip string) (string, error) {
	// Check if the passed value is a valid IP address
	if !ValidateIPAddress(ip) {
		return "", errors.New("invalid IP address")
	}
	// Split the IP it its 4 octets
	octets := strings.Split(ip, ".")

	// Replace the last octet
	octets[len(octets)-1] = "1"

	// Join the octets into an IP address
	return strings.Join(octets, "."), nil
}

// Checks if an IP address is valid
func ValidateIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

// Creates a random Mac Address
func RandomMacAddress() (string, error) {
	buf := make([]byte, 3)
	rand.Read(buf)

	buf[0] |= 2
	mac := fmt.Sprintf("52:54:00:%02x:%02x:%02x", buf[0], buf[1], buf[2])
	return mac, nil
}

// Generate a random IP address
func GenerateIPAddress() string {
	rnd.New(rnd.NewSource(0))

	min := 10
	max := 254
	octet := rnd.Intn(max-min) + min
	return fmt.Sprintf("%s.%d", ipRange, octet)

}

// Download fetches a file from the internet
func Download(url string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed to download image with error '%d'", resp.StatusCode))
	}

	return io.ReadAll(resp.Body)

}

// Download fetches a file from the internet
func DownloadAndSave(url, destination string) error {
	// Get the data
	data, err := Download(url)
	if err != nil {
		return err
	}

	// get the filename from the URl
	fileName, err := imgutil.GetFilenameFromURL(url)
	if err != nil {
		return err
	}

	// Create the file
	err = os.WriteFile(filepath.Join(destination, fileName), data, 0644)
	if err != nil {
		return err
	}

	return nil

}

func GetIPAddress(net string) string {
	ip := strings.Split(net, "/")
	return ip[0]
}
