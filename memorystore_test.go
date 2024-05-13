package main

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func TestDHCPV4Options(t *testing.T) {
	store := NewMemoryStore()

	// Create a new node with DHCP options
	node := &Node{
		MAC: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		OptionsV4: dhcpv4.Options{
			uint8(dhcpv4.OptionSubnetMask):       []byte{255, 255, 255, 0},
			uint8(dhcpv4.OptionRouter):           []byte{192, 168, 1, 1},
			uint8(dhcpv4.OptionDomainNameServer): []byte{8, 8, 8, 8},
			uint8(dhcpv4.OptionBootfileName):     []byte("pxelinux.0"),
		},
	}

	// Add the node to the store
	store.SetNode(node.MAC.String(), node)

	// Retrieve the node from the store
	retrievedNode, ok := store.GetNode(node.MAC.String())
	if !ok {
		t.Fatalf("Node not found in store")
	}

	// Check if the DHCP options are present
	if _, ok := retrievedNode.OptionsV4[uint8(dhcpv4.OptionSubnetMask)]; !ok {
		t.Errorf("DHCPv4 option not found")
	}
}
