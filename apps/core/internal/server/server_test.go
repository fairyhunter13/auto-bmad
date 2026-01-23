// Package server provides the JSON-RPC 2.0 server over stdio.
package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"strings"
	"testing"
	"time"
)

// newTestServer creates a test server with temporary project directory
func newTestServer(t *testing.T, stdin io.Reader, stdout io.Writer, logger *log.Logger) *Server {
	t.Helper()
	projectPath := t.TempDir()
	return New(stdin, stdout, logger, projectPath)
}

func TestServerHandlesSingleRequest(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
	server.RegisterHandler("test.echo", func(params json.RawMessage) (interface{}, error) {
		return "echoed", nil
	})

	// Run server in goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	// Write a request
	go func() {
		writeFrame(stdinW, `{"jsonrpc":"2.0","method":"test.echo","id":1}`)
	}()

	// Read response
	resp := readResponse(t, stdoutR)
	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC: got %q, want %q", resp.JSONRPC, "2.0")
	}
	if resp.Result != "echoed" {
		t.Errorf("Result: got %v, want %q", resp.Result, "echoed")
	}
	if resp.ID != float64(1) {
		t.Errorf("ID: got %v, want %v", resp.ID, float64(1))
	}
}

func TestServerMethodNotFound(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
	// No handlers registered

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeFrame(stdinW, `{"jsonrpc":"2.0","method":"unknown.method","id":1}`)
	}()

	resp := readResponse(t, stdoutR)
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("Error.Code: got %d, want %d", resp.Error.Code, ErrCodeMethodNotFound)
	}
}

func TestServerInvalidRequest(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "missing jsonrpc field",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdinR, stdinW := io.Pipe()
			stdoutR, stdoutW := io.Pipe()

			server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go server.Run(ctx)

			go func() {
				writeFrame(stdinW, tt.payload)
			}()

			resp := readResponse(t, stdoutR)
			if resp.Error == nil {
				t.Fatal("expected error response")
			}
			if resp.Error.Code != ErrCodeInvalidRequest {
				t.Errorf("Error.Code: got %d, want %d", resp.Error.Code, ErrCodeInvalidRequest)
			}
		})
	}
}

func TestServerParseError(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	// Write invalid JSON
	go func() {
		writeFrame(stdinW, `{not valid json}`)
	}()

	resp := readResponse(t, stdoutR)
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != ErrCodeParseError {
		t.Errorf("Error.Code: got %d, want %d", resp.Error.Code, ErrCodeParseError)
	}
}

func TestServerMultipleRequests(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
	server.RegisterHandler("add", func(params json.RawMessage) (interface{}, error) {
		return "added", nil
	})
	server.RegisterHandler("sub", func(params json.RawMessage) (interface{}, error) {
		return "subtracted", nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	// Write multiple requests
	go func() {
		writeFrame(stdinW, `{"jsonrpc":"2.0","method":"add","id":1}`)
		writeFrame(stdinW, `{"jsonrpc":"2.0","method":"sub","id":2}`)
		writeFrame(stdinW, `{"jsonrpc":"2.0","method":"add","id":3}`)
	}()

	// Read all responses
	reader := NewMessageReader(stdoutR)
	for i := 0; i < 3; i++ {
		_, err := reader.ReadRequest()
		if err != nil {
			t.Fatalf("message %d: ReadRequest failed: %v", i+1, err)
		}
	}
}

func TestServerNotificationNoResponse(t *testing.T) {
	// Use a buffer for stdout so we can check nothing was written
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	called := make(chan struct{}, 1)
	server := newTestServer(t, stdin, stdout, log.New(io.Discard, "", 0))
	server.RegisterHandler("log.info", func(params json.RawMessage) (interface{}, error) {
		called <- struct{}{}
		return nil, nil
	})

	// Write notification (no ID) then close stdin
	writeFrame(stdin, `{"jsonrpc":"2.0","method":"log.info","params":{"msg":"hello"}}`)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	server.Run(ctx)

	select {
	case <-called:
		// Good, handler was called
	default:
		t.Error("handler was not called")
	}

	// No response should be written for notifications
	if stdout.Len() > 0 {
		t.Errorf("notification should not produce a response, got %d bytes", stdout.Len())
	}
}

func TestServerPreservesRequestID(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantID  interface{}
	}{
		{
			name:    "numeric ID",
			payload: `{"jsonrpc":"2.0","method":"test","id":42}`,
			wantID:  float64(42),
		},
		{
			name:    "string ID",
			payload: `{"jsonrpc":"2.0","method":"test","id":"my-request-id"}`,
			wantID:  "my-request-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdinR, stdinW := io.Pipe()
			stdoutR, stdoutW := io.Pipe()

			server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
			server.RegisterHandler("test", func(params json.RawMessage) (interface{}, error) {
				return "ok", nil
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go server.Run(ctx)

			go func() {
				writeFrame(stdinW, tt.payload)
			}()

			resp := readResponse(t, stdoutR)
			if resp.ID != tt.wantID {
				t.Errorf("ID: got %v (%T), want %v (%T)", resp.ID, resp.ID, tt.wantID, tt.wantID)
			}
		})
	}
}

func TestServerNullIDResponse(t *testing.T) {
	// Null ID is valid per JSON-RPC 2.0, but unusual. Let's test it explicitly.
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
	server.RegisterHandler("test", func(params json.RawMessage) (interface{}, error) {
		return "ok", nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeFrame(stdinW, `{"jsonrpc":"2.0","method":"test","id":null}`)
	}()

	resp := readResponse(t, stdoutR)
	// A request with null ID is still a request (not notification), so it gets a response
	if resp.Result != "ok" {
		t.Errorf("Result: got %v, want %q", resp.Result, "ok")
	}
}

func TestServerGracefulShutdown(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	_ = stdoutR // unused in this test

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- server.Run(ctx)
	}()

	// Give server time to start
	time.Sleep(10 * time.Millisecond)

	// Signal shutdown
	cancel()
	stdinW.Close() // Also close stdin to unblock any pending reads

	// Wait for server to stop with timeout
	select {
	case err := <-done:
		// Should return without error (clean shutdown)
		if err != nil && err != context.Canceled {
			t.Errorf("unexpected error on shutdown: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("server did not shutdown gracefully within timeout")
	}
}

func TestServerEOFShutdown(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	server := newTestServer(t, stdin, stdout, log.New(io.Discard, "", 0))

	// stdin is empty, so server should return immediately with nil (EOF = clean shutdown)
	ctx := context.Background()
	err := server.Run(ctx)
	if err != nil {
		t.Errorf("expected nil error on EOF, got: %v", err)
	}
}

func TestServerLogsToStderr(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	stderr := &bytes.Buffer{}

	server := newTestServer(t, stdinR, stdoutW, log.New(stderr, "[RPC] ", 0))
	server.RegisterHandler("test", func(params json.RawMessage) (interface{}, error) {
		return "ok", nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeFrame(stdinW, `{"jsonrpc":"2.0","method":"test","id":1}`)
	}()

	// Read response to complete the cycle
	readResponse(t, stdoutR)

	// Check that something was logged
	logOutput := stderr.String()
	if !strings.Contains(logOutput, "test") {
		t.Errorf("expected log to contain method name, got: %s", logOutput)
	}
}

func TestServerHandlerError(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
	server.RegisterHandler("test.fail", func(params json.RawMessage) (interface{}, error) {
		return nil, NewError(ErrCodeInvalidParams, "invalid parameter")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeFrame(stdinW, `{"jsonrpc":"2.0","method":"test.fail","id":1}`)
	}()

	resp := readResponse(t, stdoutR)
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != ErrCodeInvalidParams {
		t.Errorf("Error.Code: got %d, want %d", resp.Error.Code, ErrCodeInvalidParams)
	}
}

// Helper functions

func writeFrame(w io.Writer, payload string) {
	frame := make([]byte, 4+len(payload)+1)
	binary.BigEndian.PutUint32(frame[:4], uint32(len(payload)))
	copy(frame[4:], payload)
	frame[len(frame)-1] = '\n'
	w.Write(frame)
}

func readResponse(t *testing.T, r io.Reader) *Response {
	t.Helper()

	// Read length prefix
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		t.Fatalf("failed to read length prefix: %v", err)
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// Read payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		t.Fatalf("failed to read payload: %v", err)
	}

	// Read trailing newline
	nl := make([]byte, 1)
	if _, err := io.ReadFull(r, nl); err != nil {
		t.Fatalf("failed to read newline: %v", err)
	}

	// Parse as Response
	var resp Response
	if err := json.Unmarshal(payload, &resp); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}
	return &resp
}
