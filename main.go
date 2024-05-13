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

	var store NodeStorage

	if *dsn != "" {
		// Use PostgreSQL store
		pgStore, err := NewPostgresStore(*dsn)
		if err != nil {
			fmt.Printf("Failed to create PostgreSQL store: %v\n", err)
			return
		}
		store = pgStore
	} else {
		// Use memory store
		store = NewMemoryStore()
	}

	// Start selected servers based on flags
	if *apiServer {
		go startAPIServer(store)
	}
	if *tftpServer {
		go startTFTPServer("/tftpboot")
	}
	if *dhcpServer {
		go startDHCPServer(store)
	}

	select {} // Block forever
}
