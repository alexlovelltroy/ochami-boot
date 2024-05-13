package main

import (
	"encoding/json"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
)

type Node struct {
	MAC       net.HardwareAddr `json:"mac" db:"mac"`
	IPv4      net.IP           `json:"ipv4,omitempty" db:"ipv4"`
	IPv6      net.IP           `json:"ipv6,omitempty" db:"ipv6"`
	OptionsV4 dhcpv4.Options   `json:"-" db:"dhcpv4_options"`
	OptionsV6 dhcpv6.Options   `json:"-" db:"dhcpv6_options"`
}

func (n *Node) SetOptionsV4(options dhcpv4.Options) {
	n.OptionsV4 = options
}

func (n *Node) SetOptionsV6(options dhcpv6.Options) {
	n.OptionsV6 = options
}

func (n *Node) GetOptionsV4() dhcpv4.Options {
	return n.OptionsV4
}

func (n *Node) GetOptionsV6() dhcpv6.Options {
	return n.OptionsV6
}

// Custom Marshal/Unmarshal to handle HardwareAddr and IP
func (n *Node) MarshalJSON() ([]byte, error) {
	type Alias Node
	return json.Marshal(&struct {
		MAC  string `json:"mac"`
		IPv4 string `json:"ipv4,omitempty"`
		IPv6 string `json:"ipv6,omitempty"`
		*Alias
	}{
		MAC:   n.MAC.String(),
		IPv4:  n.IPv4.String(),
		IPv6:  n.IPv6.String(),
		Alias: (*Alias)(n),
	})
}

func (n *Node) UnmarshalJSON(data []byte) error {
	type Alias Node
	aux := &struct {
		MAC  string `json:"mac"`
		IPv4 string `json:"ipv4,omitempty"`
		IPv6 string `json:"ipv6,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(n),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.MAC != "" {
		hwAddr, err := net.ParseMAC(aux.MAC)
		if err != nil {
			return err
		}
		n.MAC = hwAddr
	}
	n.IPv4 = net.ParseIP(aux.IPv4)
	n.IPv6 = net.ParseIP(aux.IPv6)
	return nil
}
