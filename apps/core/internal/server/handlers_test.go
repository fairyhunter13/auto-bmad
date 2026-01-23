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

func TestSystemPingHandler(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
	RegisterSystemHandlers(server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	// Send ping request
	go func() {
		writeTestFrame(stdinW, `{"jsonrpc":"2.0","method":"system.ping","id":1}`)
	}()

	// Read response
	resp := readTestResponse(t, stdoutR)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
	if resp.Result != "pong" {
		t.Errorf("Result: got %v, want %q", resp.Result, "pong")
	}
	if resp.ID != float64(1) {
		t.Errorf("ID: got %v, want %v", resp.ID, float64(1))
	}
}

func TestSystemPingPreservesID(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantID  interface{}
	}{
		{
			name:    "numeric ID",
			payload: `{"jsonrpc":"2.0","method":"system.ping","id":42}`,
			wantID:  float64(42),
		},
		{
			name:    "string ID",
			payload: `{"jsonrpc":"2.0","method":"system.ping","id":"request-001"}`,
			wantID:  "request-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdinR, stdinW := io.Pipe()
			stdoutR, stdoutW := io.Pipe()

			server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
			RegisterSystemHandlers(server)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go server.Run(ctx)

			go func() {
				writeTestFrame(stdinW, tt.payload)
			}()

			resp := readTestResponse(t, stdoutR)
			if resp.ID != tt.wantID {
				t.Errorf("ID: got %v (%T), want %v (%T)", resp.ID, resp.ID, tt.wantID, tt.wantID)
			}
		})
	}
}

func TestSystemEchoHandler(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	server := newTestServer(t, stdinR, stdoutW, log.New(io.Discard, "", 0))
	RegisterSystemHandlers(server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	// Send echo request with params
	go func() {
		writeTestFrame(stdinW, `{"jsonrpc":"2.0","method":"system.echo","params":{"message":"hello world"},"id":1}`)
	}()

	resp := readTestResponse(t, stdoutR)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	// Result should contain the echoed message
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T", resp.Result)
	}
	if result["message"] != "hello world" {
		t.Errorf("message: got %v, want %q", result["message"], "hello world")
	}
}

// Helper functions for handler tests

func writeTestFrame(w io.Writer, payload string) {
	frame := make([]byte, 4+len(payload)+1)
	binary.BigEndian.PutUint32(frame[:4], uint32(len(payload)))
	copy(frame[4:], payload)
	frame[len(frame)-1] = '\n'
	w.Write(frame)
}

func readTestResponse(t *testing.T, r io.Reader) *Response {
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
