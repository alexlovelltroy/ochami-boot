package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type IPMACRetriever struct {
	apiURL string
}

type IPMAC struct {
	IP  string `json:"ip"`
	MAC string `json:"mac"`
}

type IPMACRepository interface {
	GetIPMACs() ([]IPMAC, error)
}

func NewIPMACRetriever(apiURL string) IPMACRepository {
	return &IPMACRetriever{
		apiURL: apiURL,
	}
}

func (r *IPMACRetriever) GetIPMACs() ([]IPMAC, error) {
	resp, err := http.Get(r.apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve IP/MAC combinations: %v", err)
	}
	defer resp.Body.Close()

	var ipmacs []IPMAC
	err = json.NewDecoder(resp.Body).Decode(&ipmacs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode IP/MAC combinations: %v", err)
	}

	return ipmacs, nil
}

type RemoteNodeStorage struct {
	retriever IPMACRepository
}

func NewRemoteNodeStorage(retriever IPMACRepository) NodeStorage {
	return &RemoteNodeStorage{
		retriever: retriever,
	}
}

func (s *RemoteNodeStorage) AddNode(node *Node) {
	// Read-only implementation, do nothing
}

func (s *RemoteNodeStorage) GetNode(mac string) (*Node, bool) {
	ipmacs, err := s.retriever.GetIPMACs()
	if err != nil {
		return nil, false
	}

	for _, ipmac := range ipmacs {
		if ipmac.MAC == mac {
			ip := net.ParseIP(ipmac.IP)          // Convert string to net.IP
			hwAddr, _ := net.ParseMAC(ipmac.MAC) // Convert string to net.HardwareAddr
			return &Node{
				IPv4: ip,
				MAC:  hwAddr, // Assign the converted net.HardwareAddr value
			}, true
		}
	}

	return nil, false
}

func (s *RemoteNodeStorage) SetNode(mac string, node *Node) {
	// Read-only implementation, do nothing
}

func (s *RemoteNodeStorage) DeleteNode(mac string) {
	// Read-only implementation, do nothing
}

func example() {
	// Example usage
	retriever := NewIPMACRetriever("https://api.example.com/ipmacs")
	ipmacs, err := retriever.GetIPMACs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("IP/MAC combinations:")
	for _, ipmac := range ipmacs {
		fmt.Printf("IP: %s, MAC: %s\n", ipmac.IP, ipmac.MAC)
	}
}
