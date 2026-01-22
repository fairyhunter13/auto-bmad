// Package server provides the JSON-RPC 2.0 server over stdio.
package server

import (
	"encoding/json"
	"fmt"
)

// Standard JSON-RPC 2.0 error codes
const (
	ErrCodeParseError     = -32700 // Invalid JSON was received
	ErrCodeInvalidRequest = -32600 // JSON is not a valid Request object
	ErrCodeMethodNotFound = -32601 // Method does not exist / is not available
	ErrCodeInvalidParams  = -32602 // Invalid method parameter(s)
	ErrCodeInternalError  = -32603 // Internal JSON-RPC error
)

// Application-specific error codes (use -32000 to -32099)
const (
	ErrCodeOpenCodeNotFound = -32001
	ErrCodeGitNotFound      = -32002
	ErrCodeJourneyNotFound  = -32003
)

// Request represents a JSON-RPC 2.0 request.
// Per spec: jsonrpc and method are required, params and id are optional.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"` // string, number, or null
	hasID   bool            // true if id field was present in JSON (even if null)
}

// UnmarshalJSON implements custom JSON unmarshaling to detect presence of id field.
func (r *Request) UnmarshalJSON(data []byte) error {
	// Use an alias to avoid recursion
	type RequestAlias Request
	aux := &struct {
		*RequestAlias
		ID json.RawMessage `json:"id"`
	}{
		RequestAlias: (*RequestAlias)(r),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Check if id field was present
	if aux.ID != nil {
		r.hasID = true
		// Parse the actual ID value
		if string(aux.ID) == "null" {
			r.ID = nil
		} else {
			if err := json.Unmarshal(aux.ID, &r.ID); err != nil {
				return err
			}
		}
	} else {
		r.hasID = false
		r.ID = nil
	}

	return nil
}

// Response represents a JSON-RPC 2.0 response.
// Per spec: Either result or error MUST be included, but not both.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// Error represents a JSON-RPC 2.0 error object.
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error implements the error interface for Error.
func (e *Error) Error() string {
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}

// NewError creates a new Error with the given code and message.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithData creates a new Error with code, message, and additional data.
func NewErrorWithData(code int, message string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewSuccessResponse creates a success Response with the given result.
func NewSuccessResponse(id interface{}, result interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
}

// NewErrorResponse creates an error Response with the given error code and message.
func NewErrorResponse(id interface{}, code int, message string) *Response {
	return &Response{
		JSONRPC: "2.0",
		Error:   NewError(code, message),
		ID:      id,
	}
}

// NewErrorResponseWithData creates an error Response with additional data.
func NewErrorResponseWithData(id interface{}, code int, message string, data interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		Error:   NewErrorWithData(code, message, data),
		ID:      id,
	}
}

// IsValid checks if the request has the required fields per JSON-RPC 2.0 spec.
func (r *Request) IsValid() bool {
	return r.JSONRPC == "2.0" && r.Method != ""
}

// IsNotification returns true if this is a notification (no ID field).
// Per JSON-RPC 2.0 spec, notifications MUST NOT have an id member.
// Note: A request with "id": null is still a request, not a notification.
func (r *Request) IsNotification() bool {
	return !r.hasID
}
