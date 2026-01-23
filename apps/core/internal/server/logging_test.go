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

func TestLoggingIncludesRequestMethod(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	stderr := &bytes.Buffer{}

	server := newTestServer(t, stdinR, stdoutW, log.New(stderr, "[RPC] ", log.LstdFlags))
	RegisterSystemHandlers(server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeLogTestFrame(stdinW, `{"jsonrpc":"2.0","method":"system.ping","id":1}`)
	}()

	// Read response to complete the request
	readLogTestResponse(t, stdoutR)

	logOutput := stderr.String()
	if !strings.Contains(logOutput, "system.ping") {
		t.Errorf("log should contain method name, got: %s", logOutput)
	}
}

func TestLoggingIncludesRequestID(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	stderr := &bytes.Buffer{}

	server := newTestServer(t, stdinR, stdoutW, log.New(stderr, "[RPC] ", log.LstdFlags))
	RegisterSystemHandlers(server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeLogTestFrame(stdinW, `{"jsonrpc":"2.0","method":"system.ping","id":12345}`)
	}()

	readLogTestResponse(t, stdoutR)

	logOutput := stderr.String()
	if !strings.Contains(logOutput, "12345") {
		t.Errorf("log should contain request ID, got: %s", logOutput)
	}
}

func TestLoggingIncludesTimestamp(t *testing.T) {
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	stderr := &bytes.Buffer{}

	// Use LstdFlags which includes timestamp
	server := newTestServer(t, stdinR, stdoutW, log.New(stderr, "[RPC] ", log.LstdFlags))
	RegisterSystemHandlers(server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	go func() {
		writeLogTestFrame(stdinW, `{"jsonrpc":"2.0","method":"system.ping","id":1}`)
	}()

	readLogTestResponse(t, stdoutR)

	logOutput := stderr.String()
	// LstdFlags includes date and time in format "2009/01/23 01:23:23"
	// Check for year/month/day pattern
	if !strings.Contains(logOutput, "/") {
		t.Errorf("log should contain timestamp (date with /), got: %s", logOutput)
	}
}

func TestLoggingForErrorResponses(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	writeLogTestFrame(stdin, `{"jsonrpc":"2.0","method":"unknown","id":1}`)

	server := newTestServer(t, stdin, stdout, log.New(stderr, "[RPC] ", 0))
	// No handlers registered - will get method not found

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	server.Run(ctx)

	logOutput := stderr.String()
	// Should log both the request and the error response
	if !strings.Contains(logOutput, "unknown") {
		t.Errorf("log should contain method name, got: %s", logOutput)
	}
	// Error responses are logged with the error code
	if !strings.Contains(logOutput, "error") {
		t.Errorf("log should contain 'error' for error responses, got: %s", logOutput)
	}
}

func TestStdoutOnlyContainsJSONRPC(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	writeLogTestFrame(stdin, `{"jsonrpc":"2.0","method":"system.ping","id":1}`)

	server := newTestServer(t, stdin, stdout, log.New(stderr, "[RPC] ", 0))
	RegisterSystemHandlers(server)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	server.Run(ctx)

	// stdout should only contain the JSON-RPC response, no log messages
	stdoutBytes := stdout.Bytes()

	if len(stdoutBytes) > 5 {
		// Skip 4-byte length prefix
		length := binary.BigEndian.Uint32(stdoutBytes[:4])
		jsonContent := stdoutBytes[4 : 4+length]

		// Should be valid JSON-RPC response
		var resp Response
		if err := json.Unmarshal(jsonContent, &resp); err != nil {
			t.Errorf("stdout should contain valid JSON-RPC, got error: %v", err)
		}
		if resp.Result != "pong" {
			t.Errorf("expected pong result, got %v", resp.Result)
		}
	}

	// Log output should be in stderr, not stdout
	stdoutStr := string(stdoutBytes)
	if strings.Contains(stdoutStr, "[RPC]") {
		t.Error("log prefix should not appear in stdout")
	}

	// Verify stderr has the logs
	if !strings.Contains(stderr.String(), "system.ping") {
		t.Error("stderr should contain log messages")
	}
}

// Helper functions

func writeLogTestFrame(w io.Writer, payload string) {
	frame := make([]byte, 4+len(payload)+1)
	binary.BigEndian.PutUint32(frame[:4], uint32(len(payload)))
	copy(frame[4:], payload)
	frame[len(frame)-1] = '\n'
	w.Write(frame)
}

func readLogTestResponse(t *testing.T, r io.Reader) *Response {
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
