package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define flags
	apiServer := flag.Bool("api", false, "Start API server")
	tftpServer := flag.Bool("tftp", false, "Start TFTP server")
	dhcpServer := flag.Bool("dhcp", false, "Start DHCP server")
	dsn := flag.String("dsn", "", "PostgreSQL DSN")

	// Parse command-line flags
	flag.Parse()

	var nodeStore NodeStorage
	macStore := NewMacMemoryStore()

	if *dsn != "" {
		// Use PostgreSQL store
		pgStore, err := NewPostgresStore(*dsn)
		if err != nil {
			fmt.Printf("Failed to create PostgreSQL store: %v\n", err)
			return
		}
		nodeStore = pgStore
	} else {
		// Use memory store
		nodeStore = NewNodeMemoryStore()

	}

	// Start selected servers based on flags
	if *apiServer {
		go startAPIServer(nodeStore, macStore)
	}
	if *tftpServer {
		go startTFTPServer("/tftpboot")
	}
	if *dhcpServer {
		go startDHCPServer(nodeStore, macStore)
	}

	select {} // Block forever
}
