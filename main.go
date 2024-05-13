package main

func main() {
	store := NewMemoryStore()
	go startCRUDServer(store)
	go startTFTPServer("/tftpboot")
	go startDHCPServer(store)

	select {} // Block forever
}
