// Package server provides the JSON-RPC 2.0 server over stdio.
package server

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

// MaxMessageSize defines the maximum allowed message size (1MB)
// This prevents DoS attacks via unbounded memory allocation.
// Source: architecture.md line 375 - "1MB buffer size limit"
const MaxMessageSize = 1024 * 1024 // 1 MB

// ErrMessageTooLarge is returned when a message exceeds MaxMessageSize
var ErrMessageTooLarge = errors.New("message size exceeds maximum allowed (1MB)")

// MessageReader reads length-prefixed JSON-RPC messages from an io.Reader.
// Frame format: [4 bytes: big-endian length][N bytes: JSON payload][1 byte: newline]
type MessageReader struct {
	reader *bufio.Reader
}

// NewMessageReader creates a new MessageReader that reads from the given reader.
func NewMessageReader(r io.Reader) *MessageReader {
	return &MessageReader{reader: bufio.NewReader(r)}
}

// ReadRequest reads and parses a single JSON-RPC request from the stream.
// Returns io.EOF when no more messages are available.
func (mr *MessageReader) ReadRequest() (*Request, error) {
	// 1. Read 4-byte length prefix (big-endian uint32)
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(mr.reader, lengthBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// 2. Validate message size to prevent DoS attacks
	if length > MaxMessageSize {
		return nil, ErrMessageTooLarge
	}

	// 3. Read JSON payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(mr.reader, payload); err != nil {
		return nil, err
	}

	// 4. Read and discard trailing newline
	if _, err := mr.reader.ReadByte(); err != nil {
		return nil, err
	}

	// 5. Parse JSON into Request
	var req Request
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

// MessageWriter writes length-prefixed JSON-RPC messages to an io.Writer.
// Frame format: [4 bytes: big-endian length][N bytes: JSON payload][1 byte: newline]
type MessageWriter struct {
	writer io.Writer
}

// NewMessageWriter creates a new MessageWriter that writes to the given writer.
func NewMessageWriter(w io.Writer) *MessageWriter {
	return &MessageWriter{writer: w}
}

// WriteResponse serializes and writes a Response with length-prefixed framing.
func (mw *MessageWriter) WriteResponse(resp *Response) error {
	return mw.writeJSON(resp)
}

// WriteRequest serializes and writes a Request with length-prefixed framing.
func (mw *MessageWriter) WriteRequest(req *Request) error {
	return mw.writeJSON(req)
}

// writeJSON marshals any value to JSON and writes it with length-prefixed framing.
func (mw *MessageWriter) writeJSON(v interface{}) error {
	// 1. Marshal to JSON
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// 2. Write 4-byte length prefix (big-endian)
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, uint32(len(payload)))
	if _, err := mw.writer.Write(lengthBuf); err != nil {
		return err
	}

	// 3. Write JSON payload
	if _, err := mw.writer.Write(payload); err != nil {
		return err
	}

	// 4. Write trailing newline
	if _, err := mw.writer.Write([]byte{'\n'}); err != nil {
		return err
	}

	return nil
}
