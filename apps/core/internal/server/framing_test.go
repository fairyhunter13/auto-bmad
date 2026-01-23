// Package server provides the JSON-RPC 2.0 server over stdio.
package server

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"
)

func TestMessageWriterWriteMessage(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
	}{
		{
			name: "simple pong response",
			response: &Response{
				JSONRPC: "2.0",
				Result:  "pong",
				ID:      float64(1),
			},
		},
		{
			name: "error response",
			response: &Response{
				JSONRPC: "2.0",
				Error: &Error{
					Code:    ErrCodeMethodNotFound,
					Message: "Method not found",
				},
				ID: float64(1),
			},
		},
		{
			name: "complex result",
			response: &Response{
				JSONRPC: "2.0",
				Result: map[string]interface{}{
					"status": "ok",
					"data":   []interface{}{1, 2, 3},
				},
				ID: "abc-123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			writer := NewMessageWriter(buf)

			err := writer.WriteResponse(tt.response)
			if err != nil {
				t.Fatalf("WriteResponse failed: %v", err)
			}

			// Verify frame format: [4-byte length][JSON payload][\n]
			data := buf.Bytes()
			if len(data) < 5 {
				t.Fatalf("output too short: %d bytes", len(data))
			}

			// Read length prefix
			length := binary.BigEndian.Uint32(data[:4])

			// Check length matches payload
			expectedLen := len(data) - 4 - 1 // minus length prefix and newline
			if int(length) != expectedLen {
				t.Errorf("length prefix: got %d, want %d", length, expectedLen)
			}

			// Check trailing newline
			if data[len(data)-1] != '\n' {
				t.Errorf("missing trailing newline")
			}

			// Verify payload is valid JSON
			payload := data[4 : len(data)-1]
			if !bytes.Contains(payload, []byte(`"jsonrpc":"2.0"`)) {
				t.Errorf("payload doesn't contain jsonrpc field: %s", payload)
			}
		})
	}
}

func TestMessageReaderReadMessage(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		wantMethod string
		wantID     interface{}
		wantErr    bool
	}{
		{
			name:       "simple ping request",
			payload:    `{"jsonrpc":"2.0","method":"system.ping","id":1}`,
			wantMethod: "system.ping",
			wantID:     float64(1),
		},
		{
			name:       "request with string ID",
			payload:    `{"jsonrpc":"2.0","method":"journey.start","id":"req-001"}`,
			wantMethod: "journey.start",
			wantID:     "req-001",
		},
		{
			name:       "request with params",
			payload:    `{"jsonrpc":"2.0","method":"config.set","params":{"key":"value"},"id":42}`,
			wantMethod: "config.set",
			wantID:     float64(42),
		},
		{
			name:       "notification (no id)",
			payload:    `{"jsonrpc":"2.0","method":"log.info"}`,
			wantMethod: "log.info",
			wantID:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create framed message: [length][payload][\n]
			frame := createFrame(tt.payload)
			reader := NewMessageReader(bytes.NewReader(frame))

			req, err := reader.ReadRequest()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ReadRequest failed: %v", err)
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

func TestFramingRoundtrip(t *testing.T) {
	// Write a response, read it back as a request-like message
	buf := &bytes.Buffer{}
	writer := NewMessageWriter(buf)

	originalResp := &Response{
		JSONRPC: "2.0",
		Result:  "pong",
		ID:      float64(42),
	}

	if err := writer.WriteResponse(originalResp); err != nil {
		t.Fatalf("WriteResponse failed: %v", err)
	}

	// Read the frame back
	reader := NewMessageReader(bytes.NewReader(buf.Bytes()))
	req, err := reader.ReadRequest()
	if err != nil {
		t.Fatalf("ReadRequest failed: %v", err)
	}

	// Verify the data survived roundtrip
	if req.JSONRPC != "2.0" {
		t.Errorf("JSONRPC: got %q, want %q", req.JSONRPC, "2.0")
	}
}

func TestMessageReaderMultipleMessages(t *testing.T) {
	// Create multiple framed messages
	buf := &bytes.Buffer{}
	buf.Write(createFrame(`{"jsonrpc":"2.0","method":"first","id":1}`))
	buf.Write(createFrame(`{"jsonrpc":"2.0","method":"second","id":2}`))
	buf.Write(createFrame(`{"jsonrpc":"2.0","method":"third","id":3}`))

	reader := NewMessageReader(bytes.NewReader(buf.Bytes()))

	methods := []string{"first", "second", "third"}
	for i, wantMethod := range methods {
		req, err := reader.ReadRequest()
		if err != nil {
			t.Fatalf("message %d: ReadRequest failed: %v", i+1, err)
		}
		if req.Method != wantMethod {
			t.Errorf("message %d: Method: got %q, want %q", i+1, req.Method, wantMethod)
		}
	}

	// Fourth read should return EOF
	_, err := reader.ReadRequest()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestMessageReaderEOF(t *testing.T) {
	reader := NewMessageReader(bytes.NewReader([]byte{}))
	_, err := reader.ReadRequest()
	if err != io.EOF {
		t.Errorf("expected EOF on empty input, got %v", err)
	}
}

func TestMessageReaderPartialLength(t *testing.T) {
	// Only 2 bytes of the 4-byte length prefix
	reader := NewMessageReader(bytes.NewReader([]byte{0x00, 0x00}))
	_, err := reader.ReadRequest()
	if err == nil {
		t.Error("expected error on partial length prefix")
	}
}

func TestMessageReaderPartialPayload(t *testing.T) {
	// Full length prefix but truncated payload
	frame := make([]byte, 4)
	binary.BigEndian.PutUint32(frame, 100) // Says 100 bytes
	frame = append(frame, []byte("short")...)

	reader := NewMessageReader(bytes.NewReader(frame))
	_, err := reader.ReadRequest()
	if err == nil {
		t.Error("expected error on truncated payload")
	}
}

func TestMessageReaderInvalidJSON(t *testing.T) {
	frame := createFrame(`{invalid json}`)
	reader := NewMessageReader(bytes.NewReader(frame))
	_, err := reader.ReadRequest()
	if err == nil {
		t.Error("expected error on invalid JSON")
	}
}

func TestMessageReaderHandlesNewlinesInPayload(t *testing.T) {
	// This is the CRITICAL test - payloads with newlines must work
	payload := `{"jsonrpc":"2.0","method":"code.submit","params":{"code":"line1\nline2\nline3"},"id":1}`
	frame := createFrame(payload)

	reader := NewMessageReader(bytes.NewReader(frame))
	req, err := reader.ReadRequest()
	if err != nil {
		t.Fatalf("ReadRequest failed on payload with newlines: %v", err)
	}

	if req.Method != "code.submit" {
		t.Errorf("Method: got %q, want %q", req.Method, "code.submit")
	}

	// Verify params contain the newlines
	if !strings.Contains(string(req.Params), "line1\\nline2\\nline3") {
		t.Errorf("params don't contain expected newlines: %s", req.Params)
	}
}

func TestMessageWriterRequestRoundtrip(t *testing.T) {
	// Test writing and reading a request
	buf := &bytes.Buffer{}
	writer := NewMessageWriter(buf)

	originalReq := &Request{
		JSONRPC: "2.0",
		Method:  "test.method",
		ID:      float64(99),
	}

	if err := writer.WriteRequest(originalReq); err != nil {
		t.Fatalf("WriteRequest failed: %v", err)
	}

	reader := NewMessageReader(bytes.NewReader(buf.Bytes()))
	req, err := reader.ReadRequest()
	if err != nil {
		t.Fatalf("ReadRequest failed: %v", err)
	}

	if req.JSONRPC != originalReq.JSONRPC {
		t.Errorf("JSONRPC mismatch")
	}
	if req.Method != originalReq.Method {
		t.Errorf("Method mismatch")
	}
	if req.ID != originalReq.ID {
		t.Errorf("ID mismatch")
	}
}

// TestMessageReaderDoSProtection verifies that messages exceeding MaxMessageSize are rejected
func TestMessageReaderDoSProtection(t *testing.T) {
	tests := []struct {
		name        string
		length      uint32
		wantErr     error
		wantErrType string
	}{
		{
			name:    "valid small message",
			length:  100,
			wantErr: nil,
		},
		{
			name:        "oversized by 1 byte",
			length:      MaxMessageSize + 1,
			wantErr:     ErrMessageTooLarge,
			wantErrType: "ErrMessageTooLarge",
		},
		{
			name:        "extremely large (100MB)",
			length:      100 * 1024 * 1024,
			wantErr:     ErrMessageTooLarge,
			wantErrType: "ErrMessageTooLarge",
		},
		{
			name:        "maximum uint32 (4GB) - DoS attack vector",
			length:      0xFFFFFFFF,
			wantErr:     ErrMessageTooLarge,
			wantErrType: "ErrMessageTooLarge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create frame with specified length prefix
			frame := make([]byte, 4)
			binary.BigEndian.PutUint32(frame, tt.length)

			// If valid size, append actual payload + newline
			if tt.wantErr == nil {
				payload := `{"jsonrpc":"2.0","method":"test","id":1}`
				// Pad to exact length if needed
				if tt.length > uint32(len(payload)) {
					payload += strings.Repeat(" ", int(tt.length)-len(payload))
				} else {
					payload = payload[:tt.length]
				}
				frame = append(frame, []byte(payload)...)
				frame = append(frame, '\n')
			}
			// For error cases, don't append payload (test will fail on length check)

			reader := NewMessageReader(bytes.NewReader(frame))
			_, err := reader.ReadRequest()

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErrType)
					return
				}
				if err != tt.wantErr {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for valid size: %v", err)
				}
			}
		})
	}
}

// createFrame creates a length-prefixed frame from a JSON payload
func createFrame(payload string) []byte {
	frame := make([]byte, 4+len(payload)+1)
	binary.BigEndian.PutUint32(frame[:4], uint32(len(payload)))
	copy(frame[4:], payload)
	frame[len(frame)-1] = '\n'
	return frame
}
