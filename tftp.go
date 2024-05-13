package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pin/tftp/v3"
)

// Hander for read (aka GET) requests.
func readHandler(baseDir string) func(string, io.ReaderFrom) error {
	return func(filename string, rf io.ReaderFrom) error {
		// Ensure the requested file is within the base directory
		requestedPath := filepath.Join(baseDir, filename)
		if !strings.HasPrefix(requestedPath, baseDir) {
			return fmt.Errorf("illegal file access: %s", filename)
		}
		file, err := os.Open(requestedPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "opening %s: %v\n", filename, err)
			return err
		}
		n, err := rf.ReadFrom(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading %s: %v\n", filename, err)
			return err
		}
		fmt.Printf("%d bytes sent\n", n)
		return nil
	}
}

// Hook for logging on every transfer completion or failure.
type logHook struct{}

func (h *logHook) OnSuccess(stats tftp.TransferStats) {
	fmt.Printf("Transfer of %s to %s complete\n", stats.Filename, stats.RemoteAddr)
}
func (h *logHook) OnFailure(stats tftp.TransferStats, err error) {
	fmt.Printf("Transfer of %s to %s failed: %v\n", stats.Filename, stats.RemoteAddr, err)
}

func startTFTPServer(baseDir string) {
	// Start the server.
	s := tftp.NewServer(readHandler(baseDir), nil)
	s.SetHook(&logHook{})
	go func() {
		err := s.ListenAndServe(":69")
		if err != nil {
			fmt.Fprintf(os.Stdout, "Can't start the server: %v\n", err)
			os.Exit(1)
		}
	}()
}
