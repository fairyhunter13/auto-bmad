// Package server provides the JSON-RPC 2.0 server over stdio.
package server

import (
	"encoding/json"
)

// Build-time variables set via ldflags
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// SetVersionInfo sets the version information for the system.version handler.
// This should be called from main() with the ldflags values.
func SetVersionInfo(version, commit, date string) {
	Version = version
	Commit = commit
	Date = date
}

// RegisterSystemHandlers registers the built-in system methods.
// These are utility methods for health checking and debugging.
func RegisterSystemHandlers(s *Server) {
	s.RegisterHandler("system.ping", handleSystemPing)
	s.RegisterHandler("system.echo", handleSystemEcho)
	s.RegisterHandler("system.version", handleSystemVersion)
}

// handleSystemPing responds with "pong" for health checks.
// This is the simplest possible handler for testing connectivity.
func handleSystemPing(params json.RawMessage) (interface{}, error) {
	return "pong", nil
}

// EchoParams are the parameters for system.echo.
type EchoParams struct {
	Message string `json:"message"`
}

// EchoResult is the result of system.echo.
type EchoResult struct {
	Message string `json:"message"`
}

// handleSystemEcho echoes back the provided message.
// Useful for testing that params are being passed correctly.
func handleSystemEcho(params json.RawMessage) (interface{}, error) {
	var p EchoParams
	if params != nil {
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", err.Error())
		}
	}
	return EchoResult{Message: p.Message}, nil
}

// VersionResult is the result of system.version.
type VersionResult struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// handleSystemVersion returns the backend version information.
func handleSystemVersion(params json.RawMessage) (interface{}, error) {
	return VersionResult{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}, nil
}
