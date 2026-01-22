// Package server provides the JSON-RPC 2.0 server over stdio.
package server

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"testing"
)

// TestErrorCode32700ParseError verifies -32700 is returned for invalid JSON
func TestErrorCode32700ParseError(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := New(stdinR, stdoutW, log.New(io.Discard, "", 0))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	// Send invalid JSON
	go func() {
		writeErrorTestFrame(stdinW, `{not valid json at all}`)
	}()

	resp := readErrorTestResponse(t, stdoutR)
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != ErrCodeParseError {
		t.Errorf("expected code %d (Parse error), got %d", ErrCodeParseError, resp.Error.Code)
	}
	if resp.Error.Code != -32700 {
		t.Errorf("ErrCodeParseError should be -32700, got %d", resp.Error.Code)
	}
}

// TestErrorCode32600InvalidRequest verifies -32600 for malformed requests
func TestErrorCode32600InvalidRequest(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "missing jsonrpc version",
			payload: `{"method":"test","id":1}`,
		},
		{
			name:    "wrong jsonrpc version",
			payload: `{"jsonrpc":"1.0","method":"test","id":1}`,
		},
		{
			name:    "missing method",
			payload: `{"jsonrpc":"2.0","id":1}`,
		},
		{
			name:    "empty method",
			payload: `{"jsonrpc":"2.0","method":"","id":1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdinR, stdinW := io.Pipe()
			stdoutR, stdoutW := io.Pipe()

			server := New(stdinR, stdoutW, log.New(io.Discard, "", 0))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go server.Run(ctx)

			go func() {
				writeErrorTestFrame(stdinW, tt.payload)
			}()

			resp := readErrorTestResponse(t, stdoutR)
			if resp.Error == nil {
				t.Fatal("expected error response")
			}
			if resp.Error.Code != ErrCodeInvalidRequest {
				t.Errorf("expected code %d (Invalid Request), got %d", ErrCodeInvalidRequest, resp.Error.Code)
			}
			if resp.Error.Code != -32600 {
				t.Errorf("ErrCodeInvalidRequest should be -32600, got %d", resp.Error.Code)
			}
		})
	}
}

// TestErrorCode32601MethodNotFound verifies -32601 for unknown methods
func TestErrorCode32601MethodNotFound(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := New(stdinR, stdoutW, log.New(io.Discard, "", 0))
	// No handlers registered

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeErrorTestFrame(stdinW, `{"jsonrpc":"2.0","method":"nonexistent.method","id":1}`)
	}()

	resp := readErrorTestResponse(t, stdoutR)
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("expected code %d (Method not found), got %d", ErrCodeMethodNotFound, resp.Error.Code)
	}
	if resp.Error.Code != -32601 {
		t.Errorf("ErrCodeMethodNotFound should be -32601, got %d", resp.Error.Code)
	}
}

// TestErrorCode32602InvalidParams verifies -32602 for invalid parameters
func TestErrorCode32602InvalidParams(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := New(stdinR, stdoutW, log.New(io.Discard, "", 0))
	server.RegisterHandler("test.requiresParams", func(params json.RawMessage) (interface{}, error) {
		// This handler explicitly returns InvalidParams error
		return nil, NewError(ErrCodeInvalidParams, "missing required parameter 'name'")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeErrorTestFrame(stdinW, `{"jsonrpc":"2.0","method":"test.requiresParams","id":1}`)
	}()

	resp := readErrorTestResponse(t, stdoutR)
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != ErrCodeInvalidParams {
		t.Errorf("expected code %d (Invalid params), got %d", ErrCodeInvalidParams, resp.Error.Code)
	}
	if resp.Error.Code != -32602 {
		t.Errorf("ErrCodeInvalidParams should be -32602, got %d", resp.Error.Code)
	}
}

// TestErrorCode32603InternalError verifies -32603 for handler errors
func TestErrorCode32603InternalError(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := New(stdinR, stdoutW, log.New(io.Discard, "", 0))
	server.RegisterHandler("test.fails", func(params json.RawMessage) (interface{}, error) {
		// Return a generic error (not a JSON-RPC Error type)
		return nil, io.ErrUnexpectedEOF
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeErrorTestFrame(stdinW, `{"jsonrpc":"2.0","method":"test.fails","id":1}`)
	}()

	resp := readErrorTestResponse(t, stdoutR)
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != ErrCodeInternalError {
		t.Errorf("expected code %d (Internal error), got %d", ErrCodeInternalError, resp.Error.Code)
	}
	if resp.Error.Code != -32603 {
		t.Errorf("ErrCodeInternalError should be -32603, got %d", resp.Error.Code)
	}
}

// TestErrorResponseStructure verifies error responses have correct structure
func TestErrorResponseStructure(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := New(stdinR, stdoutW, log.New(io.Discard, "", 0))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeErrorTestFrame(stdinW, `{"jsonrpc":"2.0","method":"unknown","id":123}`)
	}()

	resp := readErrorTestResponse(t, stdoutR)

	// Per JSON-RPC 2.0 spec:
	// - jsonrpc MUST be "2.0"
	// - error MUST contain code (integer) and message (string)
	// - error MAY contain data
	// - id MUST match request id (or null if couldn't determine)
	// - result MUST NOT be present

	if resp.JSONRPC != "2.0" {
		t.Errorf("jsonrpc must be 2.0, got %q", resp.JSONRPC)
	}
	if resp.Result != nil {
		t.Error("result must not be present in error response")
	}
	if resp.Error == nil {
		t.Fatal("error must be present")
	}
	if resp.Error.Message == "" {
		t.Error("error.message must not be empty")
	}
	if resp.ID != float64(123) {
		t.Errorf("id must match request, got %v", resp.ID)
	}
}

// Helper functions

func writeErrorTestFrame(w io.Writer, payload string) {
	frame := make([]byte, 4+len(payload)+1)
	binary.BigEndian.PutUint32(frame[:4], uint32(len(payload)))
	copy(frame[4:], payload)
	frame[len(frame)-1] = '\n'
	w.Write(frame)
}

func readErrorTestResponse(t *testing.T, r io.Reader) *Response {
	t.Helper()

	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		t.Fatalf("failed to read length prefix: %v", err)
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		t.Fatalf("failed to read payload: %v", err)
	}

	nl := make([]byte, 1)
	if _, err := io.ReadFull(r, nl); err != nil {
		t.Fatalf("failed to read newline: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(payload, &resp); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}
	return &resp
}
