package server

import (
	"encoding/json"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/opencode"
)

// RegisterOpenCodeHandlers registers all OpenCode-related JSON-RPC handlers.
func RegisterOpenCodeHandlers(s *Server) {
	s.RegisterHandler("opencode.getProfiles", handleGetProfiles)
	s.RegisterHandler("opencode.detect", handleDetect)
}

// handleGetProfiles returns the list of available OpenCode profiles.
func handleGetProfiles(params json.RawMessage) (interface{}, error) {
	result, err := opencode.GetProfiles()
	if err != nil {
		return nil, NewErrorWithData(ErrCodeInternalError, "Failed to get profiles", err.Error())
	}
	return result, nil
}

// handleDetect returns OpenCode CLI detection information.
func handleDetect(params json.RawMessage) (interface{}, error) {
	result, err := opencode.Detect()
	if err != nil {
		return nil, NewErrorWithData(ErrCodeInternalError, "Failed to detect OpenCode", err.Error())
	}
	return result, nil
}
