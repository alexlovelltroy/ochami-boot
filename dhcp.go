package main

import (
	"log"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/server6"
)

func startDHCPServer(store NodeStorage) {
	dhcpv4Handler := func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		clientMAC := m.ClientHWAddr.String()
		node, exists := store.GetNode(clientMAC)
		if !exists || (node.IPv4 == nil && node.IPv6 == nil) {
			return // Do not respond if no node info or no IP specified
		}

		if node.IPv4 != nil {
			leaseTime := uint32(24 * time.Hour / time.Second) // Convert duration to uint32
			response, err := dhcpv4.NewReplyFromRequest(m,
				dhcpv4.WithServerIP(net.IP{192, 168, 1, 1}), // Your DHCP server IP
				dhcpv4.WithYourIP(node.IPv4),                // IP to be assigned
				dhcpv4.WithLeaseTime(leaseTime),             // Use the converted leaseTime value
				dhcpv4.WithMessageType(dhcpv4.MessageTypeAck),
				dhcpv4.WithOption(dhcpv4.OptServerIdentifier(net.IP{192, 168, 1, 1})),
			)
			if err != nil {
				log.Printf("Failed to create DHCPv4 ACK: %v", err)
				return
			}
			conn.WriteTo(response.ToBytes(), peer)
		}
	}

	dhcpv6Handler := func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
		if msg, ok := m.(*dhcpv6.Message); ok {
			clientMAC := msg.Options.ClientID().ToBytes()
			node, exists := store.GetNode(string(clientMAC))
			if !exists || (node.IPv6 == nil) {
				return // Do not respond if no node info or no IPv6 specified
			}

			v6addr := dhcpv6.OptIAAddress{
				IPv6Addr: node.IPv6,
			}

			response, err := dhcpv6.NewReplyFromMessage(msg,
				dhcpv6.WithServerID(msg.Options.ServerID()),
				dhcpv6.WithIANA(v6addr),
			)

			if err != nil {
				log.Printf("Failed to create DHCPv6 Reply: %v", err)
				return
			}
			conn.WriteTo(response.ToBytes(), peer)
		}
	}

	// Start DHCPv4 server
	laddr4 := net.UDPAddr{IP: net.IPv4zero, Port: dhcpv4.ServerPort}
	server4, err := server4.NewServer("0.0.0.0", &laddr4, dhcpv4Handler)
	if err != nil {
		log.Fatalf("Failed to start DHCPv4 server: %v", err)
	}
	go server4.Serve()

	laddr6 := net.UDPAddr{IP: net.IPv6zero, Port: dhcpv6.DefaultServerPort}
	server6, err := server6.NewServer("::", &laddr6, dhcpv6Handler)
	if err != nil {
		log.Fatalf("Failed to start DHCPv6 server: %v", err)
	}
	go server6.Serve()
}
