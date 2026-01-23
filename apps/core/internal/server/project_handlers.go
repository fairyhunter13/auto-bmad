package server

import (
	"encoding/json"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/checkpoint"
	"github.com/fairyhunter13/auto-bmad/apps/core/internal/opencode"
	"github.com/fairyhunter13/auto-bmad/apps/core/internal/project"
)

// RegisterProjectHandlers registers project-related JSON-RPC methods.
func RegisterProjectHandlers(s *Server) {
	s.RegisterHandler("project.detectDependencies", handleDetectDependencies)
	s.RegisterHandler("project.scan", handleProjectScan)
	s.RegisterHandler("project.getRecent", handleGetRecent)
	s.RegisterHandler("project.addRecent", handleAddRecent)
	s.RegisterHandler("project.removeRecent", handleRemoveRecent)
	s.RegisterHandler("project.setContext", handleSetContext)
}

// handleDetectDependencies detects and validates system dependencies.
// Returns detection results for OpenCode CLI and Git.
func handleDetectDependencies(params json.RawMessage) (interface{}, error) {
	opencodeResult, _ := opencode.Detect()
	gitResult, _ := checkpoint.DetectGit()

	return map[string]interface{}{
		"opencode": opencodeResult,
		"git":      gitResult,
	}, nil
}

// ScanParams represents the parameters for project.scan method
type ScanParams struct {
	Path string `json:"path"`
}

// handleProjectScan scans a project directory for BMAD structure
func handleProjectScan(params json.RawMessage) (interface{}, error) {
	var p ScanParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", err.Error())
	}

	if p.Path == "" {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", "path is required")
	}

	// Validate and sanitize the path to prevent path traversal attacks
	validatedPath, err := ValidateProjectPath(p.Path)
	if err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid path", err.Error())
	}

	result, err := project.Scan(validatedPath)
	if err != nil {
		return nil, NewErrorWithData(ErrCodeInternalError, "Failed to scan project", err.Error())
	}

	return result, nil
}

// handleGetRecent returns the list of recent projects
func handleGetRecent(params json.RawMessage) (interface{}, error) {
	rm := project.GetRecentManager()
	projects, err := rm.GetAll()
	if err != nil {
		return nil, NewErrorWithData(ErrCodeInternalError, "Failed to get recent projects", err.Error())
	}

	return projects, nil
}

// AddRecentParams represents the parameters for project.addRecent
type AddRecentParams struct {
	Path string `json:"path"`
}

// handleAddRecent adds a project to the recent list
func handleAddRecent(params json.RawMessage) (interface{}, error) {
	var p AddRecentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", err.Error())
	}

	if p.Path == "" {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", "path is required")
	}

	// Validate and sanitize the path to prevent path traversal attacks
	validatedPath, err := ValidateProjectPath(p.Path)
	if err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid path", err.Error())
	}

	rm := project.GetRecentManager()
	if err := rm.Add(validatedPath); err != nil {
		return nil, NewErrorWithData(ErrCodeInternalError, "Failed to add recent project", err.Error())
	}

	return nil, nil
}

// RemoveRecentParams represents the parameters for project.removeRecent
type RemoveRecentParams struct {
	Path string `json:"path"`
}

// handleRemoveRecent removes a project from the recent list
func handleRemoveRecent(params json.RawMessage) (interface{}, error) {
	var p RemoveRecentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", err.Error())
	}

	if p.Path == "" {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", "path is required")
	}

	// Validate path (allow non-existent since we're removing it anyway)
	validator := &PathValidator{
		AllowNonExistent: true,
		RequireDirectory: false,
	}
	validatedPath, err := validator.Validate(p.Path)
	if err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid path", err.Error())
	}

	rm := project.GetRecentManager()
	if err := rm.Remove(validatedPath); err != nil {
		return nil, NewErrorWithData(ErrCodeInternalError, "Failed to remove recent project", err.Error())
	}

	return nil, nil
}

// SetContextParams represents the parameters for project.setContext
type SetContextParams struct {
	Path    string `json:"path"`
	Context string `json:"context"`
}

// handleSetContext sets the context description for a project
func handleSetContext(params json.RawMessage) (interface{}, error) {
	var p SetContextParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", err.Error())
	}

	if p.Path == "" {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid params", "path is required")
	}

	// Validate and sanitize the path to prevent path traversal attacks
	validatedPath, err := ValidateProjectPath(p.Path)
	if err != nil {
		return nil, NewErrorWithData(ErrCodeInvalidParams, "Invalid path", err.Error())
	}

	rm := project.GetRecentManager()
	if err := rm.SetContext(validatedPath, p.Context); err != nil {
		return nil, NewErrorWithData(ErrCodeInternalError, "Failed to set project context", err.Error())
	}

	return nil, nil
}
