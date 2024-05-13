package main

func main() {
	store := NewMemoryStore()
	go startAPIServer(store)
	go startTFTPServer("/tftpboot")
	go startDHCPServer(store)

	select {} // Block forever
}
