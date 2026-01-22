// Package server provides network status handlers.
package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/network"
)

// Global network monitor instance
var networkMonitor *network.DebouncedMonitor
var networkServer *Server

// InitNetworkMonitor initializes the global network monitor.
// This should be called once during server startup.
func InitNetworkMonitor(ctx context.Context, s *Server) {
	networkServer = s

	networkMonitor = network.NewDebouncedMonitor(
		30*time.Second, // Check interval
		5*time.Second,  // Debounce window
		func(old, new network.Status) {
			// Emit status change event to frontend
			if networkServer != nil {
				event := map[string]interface{}{
					"previous": old,
					"current":  new,
				}
				if err := networkServer.EmitEvent("network.statusChanged", event); err != nil {
					// Log error but don't fail
					networkServer.logger.Printf("Failed to emit network.statusChanged event: %v", err)
				}
			}
		},
	)

	// Start monitoring in background
	go networkMonitor.Start(ctx)
}

// RegisterNetworkHandlers registers network-related JSON-RPC handlers.
func RegisterNetworkHandlers(s *Server) {
	s.RegisterHandler("network.getStatus", handleNetworkGetStatus)
}

// handleNetworkGetStatus returns the current network connectivity status.
func handleNetworkGetStatus(params json.RawMessage) (interface{}, error) {
	if networkMonitor == nil {
		// Return checking status if monitor not initialized
		return network.NetworkStatus{
			Status:      network.StatusChecking,
			LastChecked: time.Now(),
		}, nil
	}

	return networkMonitor.GetStatus(), nil
}
