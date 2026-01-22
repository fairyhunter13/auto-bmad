// Package main is the entry point for the AutoBMAD core backend.
// This binary provides JSON-RPC 2.0 server functionality for journey orchestration.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/server"
)

// Version information (set at build time)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Create logger that writes to stderr (stdout is reserved for JSON-RPC)
	logger := log.New(os.Stderr, "[RPC] ", log.LstdFlags)

	// Set version info for system.version handler
	server.SetVersionInfo(version, commit, date)

	// Print version info to stderr
	logger.Printf("AutoBMAD Core v%s (commit: %s, built: %s)", version, commit, date)
	logger.Println("Starting JSON-RPC server on stdio...")

	// Create server
	srv := server.New(os.Stdin, os.Stdout, logger)

	// Register system handlers
	server.RegisterSystemHandlers(srv)

	// Create context that cancels on SIGTERM/SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		logger.Printf("Received signal %v, shutting down...", sig)
		cancel()
	}()

	// Run server
	if err := srv.Run(ctx); err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}

	logger.Println("Server shutdown complete")
}
