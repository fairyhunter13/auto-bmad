// Package server provides Git repository status handlers.
package server

import (
	"encoding/json"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/checkpoint"
)

// RegisterGitHandlers registers Git-related JSON-RPC handlers.
func RegisterGitHandlers(s *Server) {
	s.RegisterHandler("git.getRepoStatus", handleGetRepoStatus)
}

// getRepoStatusParams represents parameters for git.getRepoStatus
type getRepoStatusParams struct {
	Path string `json:"path"`
}

// handleGetRepoStatus returns the Git repository status for a given path.
func handleGetRepoStatus(params json.RawMessage) (interface{}, error) {
	var req getRepoStatusParams
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "invalid params", err.Error())
	}

	if req.Path == "" {
		return nil, NewError(ErrCodeInvalidParams, "path is required")
	}

	status := checkpoint.GetRepoStatus(req.Path)
	return status, nil
}
