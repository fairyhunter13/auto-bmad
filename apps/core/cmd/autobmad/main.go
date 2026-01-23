// Package main is the entry point for the AutoBMAD core backend.
// This binary provides JSON-RPC 2.0 server functionality for journey orchestration.
package main

import (
	"context"
	"flag"
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
	// Parse command-line flags
	projectPath := flag.String("project-path", "", "Path to BMAD project root (required)")
	flag.Parse()

	// Validate required project path
	if *projectPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --project-path flag is required")
		fmt.Fprintln(os.Stderr, "Usage: autobmad --project-path /path/to/project")
		os.Exit(1)
	}

	// Create logger that writes to stderr (stdout is reserved for JSON-RPC)
	logger := log.New(os.Stderr, "[RPC] ", log.LstdFlags)

	// Set version info for system.version handler
	server.SetVersionInfo(version, commit, date)

	// Print version info to stderr
	logger.Printf("AutoBMAD Core v%s (commit: %s, built: %s)", version, commit, date)
	logger.Printf("Project path: %s", *projectPath)
	logger.Println("Starting JSON-RPC server on stdio...")

	// Create server with project path
	srv := server.New(os.Stdin, os.Stdout, logger, *projectPath)

	// Register system handlers
	server.RegisterSystemHandlers(srv)

	// Register project handlers
	server.RegisterProjectHandlers(srv)

	// Register OpenCode handlers
	server.RegisterOpenCodeHandlers(srv)

	// Register settings handlers (settings are now project-local)
	if err := server.RegisterSettingsHandlers(srv, *projectPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register settings handlers: %v\n", err)
		os.Exit(1)
	}

	// Create context that cancels on SIGTERM/SIGINT
	ctx, cancel := context.WithCancel(context.Background())

	// Register network handlers
	server.RegisterNetworkHandlers(srv)

	// Initialize network monitor (must be after handler registration)
	server.InitNetworkMonitor(ctx, srv)
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
