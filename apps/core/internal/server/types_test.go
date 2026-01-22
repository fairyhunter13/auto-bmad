// Package server provides the JSON-RPC 2.0 server over stdio.
package server

import (
	"encoding/json"
	"testing"
)

func TestRequestJSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		request  Request
		expected string
	}{
		{
			name: "basic request with string ID",
			request: Request{
				JSONRPC: "2.0",
				Method:  "system.ping",
				ID:      "abc123",
			},
			expected: `{"jsonrpc":"2.0","method":"system.ping","id":"abc123"}`,
		},
		{
			name: "request with numeric ID",
			request: Request{
				JSONRPC: "2.0",
				Method:  "journey.start",
				ID:      float64(1),
			},
			expected: `{"jsonrpc":"2.0","method":"journey.start","id":1}`,
		},
		{
			name: "request with params",
			request: Request{
				JSONRPC: "2.0",
				Method:  "config.set",
				Params:  json.RawMessage(`{"key":"value"}`),
				ID:      float64(42),
			},
			expected: `{"jsonrpc":"2.0","method":"config.set","params":{"key":"value"},"id":42}`,
		},
		{
			name: "notification (no ID)",
			request: Request{
				JSONRPC: "2.0",
				Method:  "log.info",
				Params:  json.RawMessage(`{"msg":"hello"}`),
			},
			expected: `{"jsonrpc":"2.0","method":"log.info","params":{"msg":"hello"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("got %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestRequestJSONDeserialization(t *testing.T) {
	tests := []struct {
		name       string
		jsonData   string
		wantMethod string
		wantID     interface{}
	}{
		{
			name:       "string ID",
			jsonData:   `{"jsonrpc":"2.0","method":"test","id":"str-id"}`,
			wantMethod: "test",
			wantID:     "str-id",
		},
		{
			name:       "numeric ID",
			jsonData:   `{"jsonrpc":"2.0","method":"test","id":123}`,
			wantMethod: "test",
			wantID:     float64(123),
		},
		{
			name:       "null ID",
			jsonData:   `{"jsonrpc":"2.0","method":"test","id":null}`,
			wantMethod: "test",
			wantID:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req Request
			if err := json.Unmarshal([]byte(tt.jsonData), &req); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if req.Method != tt.wantMethod {
				t.Errorf("Method: got %q, want %q", req.Method, tt.wantMethod)
			}
			if req.ID != tt.wantID {
				t.Errorf("ID: got %v (%T), want %v (%T)", req.ID, req.ID, tt.wantID, tt.wantID)
			}
		})
	}
}

func TestResponseJSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		response Response
		expected string
	}{
		{
			name: "success response with string result",
			response: Response{
				JSONRPC: "2.0",
				Result:  "pong",
				ID:      float64(1),
			},
			expected: `{"jsonrpc":"2.0","result":"pong","id":1}`,
		},
		{
			name: "success response with object result",
			response: Response{
				JSONRPC: "2.0",
				Result:  map[string]interface{}{"status": "ok"},
				ID:      "req-1",
			},
			expected: `{"jsonrpc":"2.0","result":{"status":"ok"},"id":"req-1"}`,
		},
		{
			name: "error response",
			response: Response{
				JSONRPC: "2.0",
				Error: &Error{
					Code:    ErrCodeMethodNotFound,
					Message: "Method not found",
					Data:    "unknown.method",
				},
				ID: float64(1),
			},
			expected: `{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found","data":"unknown.method"},"id":1}`,
		},
		{
			name: "error response without data",
			response: Response{
				JSONRPC: "2.0",
				Error: &Error{
					Code:    ErrCodeParseError,
					Message: "Parse error",
				},
				ID: nil,
			},
			expected: `{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("failed to marshal response: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("got %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Verify JSON-RPC 2.0 standard error codes
	if ErrCodeParseError != -32700 {
		t.Errorf("ErrCodeParseError: got %d, want -32700", ErrCodeParseError)
	}
	if ErrCodeInvalidRequest != -32600 {
		t.Errorf("ErrCodeInvalidRequest: got %d, want -32600", ErrCodeInvalidRequest)
	}
	if ErrCodeMethodNotFound != -32601 {
		t.Errorf("ErrCodeMethodNotFound: got %d, want -32601", ErrCodeMethodNotFound)
	}
	if ErrCodeInvalidParams != -32602 {
		t.Errorf("ErrCodeInvalidParams: got %d, want -32602", ErrCodeInvalidParams)
	}
	if ErrCodeInternalError != -32603 {
		t.Errorf("ErrCodeInternalError: got %d, want -32603", ErrCodeInternalError)
	}
}

func TestNewError(t *testing.T) {
	err := NewError(ErrCodeMethodNotFound, "Method not found")
	if err.Code != ErrCodeMethodNotFound {
		t.Errorf("Code: got %d, want %d", err.Code, ErrCodeMethodNotFound)
	}
	if err.Message != "Method not found" {
		t.Errorf("Message: got %q, want %q", err.Message, "Method not found")
	}
	if err.Data != nil {
		t.Errorf("Data: got %v, want nil", err.Data)
	}
}

func TestNewErrorWithData(t *testing.T) {
	err := NewErrorWithData(ErrCodeInvalidParams, "Invalid params", map[string]string{"field": "name"})
	if err.Code != ErrCodeInvalidParams {
		t.Errorf("Code: got %d, want %d", err.Code, ErrCodeInvalidParams)
	}
	if err.Data == nil {
		t.Error("Data should not be nil")
	}
}

func TestErrorImplementsError(t *testing.T) {
	err := NewError(ErrCodeInternalError, "something went wrong")
	var _ error = err // Compile-time check that Error implements error interface

	errMsg := err.Error()
	if errMsg != "JSON-RPC error -32603: something went wrong" {
		t.Errorf("Error(): got %q", errMsg)
	}
}

func TestNewSuccessResponse(t *testing.T) {
	resp := NewSuccessResponse("test-id", "result-value")
	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC: got %q, want %q", resp.JSONRPC, "2.0")
	}
	if resp.Result != "result-value" {
		t.Errorf("Result: got %v, want %v", resp.Result, "result-value")
	}
	if resp.Error != nil {
		t.Errorf("Error: got %v, want nil", resp.Error)
	}
	if resp.ID != "test-id" {
		t.Errorf("ID: got %v, want %v", resp.ID, "test-id")
	}
}

func TestNewErrorResponse(t *testing.T) {
	resp := NewErrorResponse(float64(42), ErrCodeMethodNotFound, "Method not found")
	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC: got %q, want %q", resp.JSONRPC, "2.0")
	}
	if resp.Result != nil {
		t.Errorf("Result: got %v, want nil", resp.Result)
	}
	if resp.Error == nil {
		t.Fatal("Error should not be nil")
	}
	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("Error.Code: got %d, want %d", resp.Error.Code, ErrCodeMethodNotFound)
	}
	if resp.ID != float64(42) {
		t.Errorf("ID: got %v, want %v", resp.ID, float64(42))
	}
}

func TestIsValidRequest(t *testing.T) {
	tests := []struct {
		name    string
		request Request
		valid   bool
	}{
		{
			name:    "valid request with ID",
			request: Request{JSONRPC: "2.0", Method: "test", ID: float64(1)},
			valid:   true,
		},
		{
			name:    "valid notification (no ID)",
			request: Request{JSONRPC: "2.0", Method: "log"},
			valid:   true,
		},
		{
			name:    "invalid - wrong jsonrpc version",
			request: Request{JSONRPC: "1.0", Method: "test", ID: float64(1)},
			valid:   false,
		},
		{
			name:    "invalid - missing jsonrpc",
			request: Request{Method: "test", ID: float64(1)},
			valid:   false,
		},
		{
			name:    "invalid - missing method",
			request: Request{JSONRPC: "2.0", ID: float64(1)},
			valid:   false,
		},
		{
			name:    "invalid - empty method",
			request: Request{JSONRPC: "2.0", Method: "", ID: float64(1)},
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.IsValid()
			if got != tt.valid {
				t.Errorf("IsValid(): got %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestIsNotification(t *testing.T) {
	// Test with JSON unmarshaling since hasID is set during unmarshal
	tests := []struct {
		name      string
		json      string
		wantNotif bool
	}{
		{
			name:      "request with numeric ID",
			json:      `{"jsonrpc":"2.0","method":"test","id":1}`,
			wantNotif: false,
		},
		{
			name:      "request with string ID",
			json:      `{"jsonrpc":"2.0","method":"test","id":"abc"}`,
			wantNotif: false,
		},
		{
			name:      "request with null ID (still a request)",
			json:      `{"jsonrpc":"2.0","method":"test","id":null}`,
			wantNotif: false,
		},
		{
			name:      "notification (no id field)",
			json:      `{"jsonrpc":"2.0","method":"log"}`,
			wantNotif: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req Request
			if err := json.Unmarshal([]byte(tt.json), &req); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if req.IsNotification() != tt.wantNotif {
				t.Errorf("IsNotification(): got %v, want %v", req.IsNotification(), tt.wantNotif)
			}
		})
	}
}
