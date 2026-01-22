# Story 1.2: Implement JSON-RPC Server Foundation

Status: review

## Story

As a **developer building the backend**,
I want **a JSON-RPC 2.0 server that reads from stdin and writes to stdout**,
So that **the Electron main process can communicate with the Golang backend**.

## Acceptance Criteria

1. **Given** the Golang binary is running  
   **When** a valid JSON-RPC request is sent via stdin  
   **Then** a valid JSON-RPC response is returned via stdout  
   **And** the response follows JSON-RPC 2.0 specification exactly

2. **Given** an invalid JSON-RPC request is sent  
   **When** the request is processed  
   **Then** proper error response with code -32600 (Invalid Request) is returned

3. **Given** a `system.ping` method is called  
   **When** the request is processed  
   **Then** the response contains `{"result": "pong"}` with matching request ID

4. **Given** an unknown method is called  
   **When** the request is processed  
   **Then** the response contains error code -32601 (Method not found)

5. **Given** the server is running  
   **When** messages contain newlines (e.g., code blocks)  
   **Then** length-prefixed framing correctly handles the payload

## Tasks / Subtasks

- [x] **Task 1: Create JSON-RPC 2.0 message types** (AC: #1, #2)
  - [x] Define `Request` struct with jsonrpc, method, params, id fields
  - [x] Define `Response` struct with jsonrpc, result, error, id fields
  - [x] Define `Error` struct with code, message, data fields
  - [x] Add JSON tags with camelCase naming

- [x] **Task 2: Implement length-prefixed framing** (AC: #5)
  - [x] Create `MessageReader` with 4-byte length prefix parsing
  - [x] Create `MessageWriter` with 4-byte length prefix encoding
  - [x] Add newline terminator after each frame
  - [x] Handle partial reads/writes correctly

- [x] **Task 3: Create JSON-RPC server core** (AC: #1, #2, #4)
  - [x] Implement message loop reading from stdin
  - [x] Parse and validate JSON-RPC 2.0 requests
  - [x] Route methods to handlers
  - [x] Write responses to stdout

- [x] **Task 4: Implement system.ping method** (AC: #3)
  - [x] Create handler returning `{"result": "pong"}`
  - [x] Ensure ID is preserved in response
  - [x] Add test coverage

- [x] **Task 5: Implement error handling** (AC: #2, #4)
  - [x] -32700 Parse error (invalid JSON)
  - [x] -32600 Invalid Request (missing required fields)
  - [x] -32601 Method not found
  - [x] -32602 Invalid params
  - [x] -32603 Internal error

- [x] **Task 6: Add structured logging** (AC: all)
  - [x] Log all incoming requests to stderr
  - [x] Log all outgoing responses to stderr
  - [x] Include timestamps and request IDs

## Dev Notes

### Architecture Requirements

**Source:** [architecture.md#Communication Architecture]

| Aspect | Decision |
|--------|----------|
| Protocol | JSON-RPC 2.0 over stdin/stdout |
| Direction | Bidirectional (request/response + event streaming) |
| Serialization | JSON |
| Error Handling | JSON-RPC error codes + custom application codes |
| JSON Field Naming | camelCase throughout IPC boundary |

### Message Framing Protocol (CRITICAL)

**Source:** [architecture.md#IPC Protocol Resilience]

**Problem:** Newline-delimited JSON breaks when messages contain newlines (e.g., code blocks).

**Solution:** Length-prefixed framing with JSON envelope.

```
┌────────────────────────────────────────────────────────────────┐
│                     Message Frame Format                        │
├────────────────────────────────────────────────────────────────┤
│  [4 bytes: length]  [N bytes: JSON payload]  [1 byte: newline] │
│                                                                 │
│  Example:                                                       │
│  \x00\x00\x00\x2F{"jsonrpc":"2.0","method":"ping","id":1}\n    │
│  └─── 47 ────────┘└──────────── 47 bytes ───────────────┘      │
└────────────────────────────────────────────────────────────────┘
```

The 4-byte length is **big-endian uint32** representing the JSON payload length (NOT including the newline).

### JSON-RPC 2.0 Types

```go
// internal/server/types.go

// Request represents a JSON-RPC 2.0 request
type Request struct {
    JSONRPC string          `json:"jsonrpc"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params,omitempty"`
    ID      interface{}     `json:"id,omitempty"` // string, number, or null
}

// Response represents a JSON-RPC 2.0 response
type Response struct {
    JSONRPC string      `json:"jsonrpc"`
    Result  interface{} `json:"result,omitempty"`
    Error   *Error      `json:"error,omitempty"`
    ID      interface{} `json:"id"`
}

// Error represents a JSON-RPC 2.0 error
type Error struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC 2.0 error codes
const (
    ErrCodeParseError     = -32700
    ErrCodeInvalidRequest = -32600
    ErrCodeMethodNotFound = -32601
    ErrCodeInvalidParams  = -32602
    ErrCodeInternalError  = -32603
)

// Application-specific error codes (use -32000 to -32099)
const (
    ErrCodeOpenCodeNotFound = -32001
    ErrCodeGitNotFound      = -32002
    ErrCodeJourneyNotFound  = -32003
)
```

### Message Reader Implementation

```go
// internal/server/framing.go

import (
    "bufio"
    "encoding/binary"
    "encoding/json"
    "io"
)

type MessageReader struct {
    reader *bufio.Reader
}

func NewMessageReader(r io.Reader) *MessageReader {
    return &MessageReader{reader: bufio.NewReader(r)}
}

func (mr *MessageReader) ReadMessage() (*Request, error) {
    // 1. Read 4-byte length prefix (big-endian)
    lengthBuf := make([]byte, 4)
    if _, err := io.ReadFull(mr.reader, lengthBuf); err != nil {
        return nil, err
    }
    length := binary.BigEndian.Uint32(lengthBuf)
    
    // 2. Read JSON payload
    payload := make([]byte, length)
    if _, err := io.ReadFull(mr.reader, payload); err != nil {
        return nil, err
    }
    
    // 3. Read and discard trailing newline
    if _, err := mr.reader.ReadByte(); err != nil {
        return nil, err
    }
    
    // 4. Parse JSON
    var req Request
    if err := json.Unmarshal(payload, &req); err != nil {
        return nil, err
    }
    
    return &req, nil
}
```

### Message Writer Implementation

```go
// internal/server/framing.go

type MessageWriter struct {
    writer io.Writer
}

func NewMessageWriter(w io.Writer) *MessageWriter {
    return &MessageWriter{writer: w}
}

func (mw *MessageWriter) WriteMessage(resp *Response) error {
    // 1. Marshal to JSON
    payload, err := json.Marshal(resp)
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
```

### Server Core Implementation

```go
// internal/server/server.go

type Server struct {
    reader   *MessageReader
    writer   *MessageWriter
    handlers map[string]Handler
    logger   *log.Logger
}

type Handler func(params json.RawMessage) (interface{}, error)

func NewServer(stdin io.Reader, stdout io.Writer) *Server {
    return &Server{
        reader:   NewMessageReader(stdin),
        writer:   NewMessageWriter(stdout),
        handlers: make(map[string]Handler),
        logger:   log.New(os.Stderr, "[RPC] ", log.LstdFlags),
    }
}

func (s *Server) RegisterHandler(method string, handler Handler) {
    s.handlers[method] = handler
}

func (s *Server) Run() error {
    for {
        req, err := s.reader.ReadMessage()
        if err == io.EOF {
            return nil // Clean shutdown
        }
        if err != nil {
            s.writeError(nil, ErrCodeParseError, "Parse error", err.Error())
            continue
        }
        
        s.logger.Printf("Request: %s (id=%v)", req.Method, req.ID)
        
        // Validate JSON-RPC version
        if req.JSONRPC != "2.0" {
            s.writeError(req.ID, ErrCodeInvalidRequest, "Invalid Request", "jsonrpc must be 2.0")
            continue
        }
        
        // Find handler
        handler, ok := s.handlers[req.Method]
        if !ok {
            s.writeError(req.ID, ErrCodeMethodNotFound, "Method not found", req.Method)
            continue
        }
        
        // Execute handler
        result, err := handler(req.Params)
        if err != nil {
            s.writeError(req.ID, ErrCodeInternalError, "Internal error", err.Error())
            continue
        }
        
        // Write success response
        s.writeResult(req.ID, result)
    }
}
```

### File Structure

```
apps/core/
├── cmd/autobmad/
│   └── main.go           # Entry point, creates and runs server
└── internal/
    └── server/
        ├── server.go     # Server struct and Run loop
        ├── framing.go    # MessageReader, MessageWriter
        ├── types.go      # Request, Response, Error structs
        ├── handlers.go   # Handler registration
        └── server_test.go # Unit tests
```

### Testing Requirements

Create tests that verify:

1. **Framing roundtrip**: Write message → Read message produces identical data
2. **system.ping**: Returns `{"jsonrpc":"2.0","result":"pong","id":1}`
3. **Method not found**: Returns error code -32601
4. **Invalid request**: Returns error code -32600
5. **Parse error**: Returns error code -32700

Example test:

```go
// internal/server/server_test.go

func TestSystemPing(t *testing.T) {
    stdin := bytes.NewBuffer(nil)
    stdout := bytes.NewBuffer(nil)
    
    server := NewServer(stdin, stdout)
    server.RegisterHandler("system.ping", func(params json.RawMessage) (interface{}, error) {
        return "pong", nil
    })
    
    // Write request
    req := &Request{JSONRPC: "2.0", Method: "system.ping", ID: 1}
    writer := NewMessageWriter(stdin)
    writer.WriteMessage(&Response{JSONRPC: "2.0", Result: "pong", ID: 1}) // Simulate
    
    // ... verify response
}
```

### Logging to stderr (IMPORTANT)

All logs MUST go to stderr, NOT stdout. Stdout is reserved for JSON-RPC messages only.

```go
// CORRECT
log.New(os.Stderr, "[RPC] ", log.LstdFlags)

// WRONG - will corrupt IPC
fmt.Println("Debug message")
log.Println("Info")
```

### Buffer Sizes

**Source:** [architecture.md#Buffer Size Specifications]

| Buffer | Size | Rationale |
|--------|------|-----------|
| Golang stdout pipe | 64KB (OS default) | Sufficient for JSON-RPC |
| IPC read buffer | 1MB | Handle large payloads (code artifacts) |

### Dependencies (Story 1.1)

This story depends on Story 1.1 completing the monorepo setup. The Golang module must exist at `apps/core/`.

### References

- [architecture.md#Communication Architecture] - IPC protocol decisions
- [architecture.md#IPC Protocol Resilience] - Message framing
- [architecture.md#Buffer Size Specifications] - Buffer sizes
- [architecture.md#Corrupted Message Recovery] - Error handling
- [prd.md#NFR-P4] - OpenCode process spawn time < 2 seconds

## File List

### New Files
- `apps/core/internal/server/types.go` - JSON-RPC 2.0 message types (Request, Response, Error)
- `apps/core/internal/server/types_test.go` - Tests for message types
- `apps/core/internal/server/framing.go` - Length-prefixed message framing (MessageReader, MessageWriter)
- `apps/core/internal/server/framing_test.go` - Tests for framing
- `apps/core/internal/server/handlers.go` - System handlers (system.ping, system.echo)
- `apps/core/internal/server/handlers_test.go` - Tests for handlers
- `apps/core/internal/server/errors_test.go` - Tests for all JSON-RPC 2.0 error codes
- `apps/core/internal/server/logging_test.go` - Tests for logging behavior
- `apps/core/cmd/autobmad/main_test.go` - Integration tests for binary

### Modified Files
- `apps/core/internal/server/server.go` - Full JSON-RPC 2.0 server implementation
- `apps/core/internal/server/server_test.go` - Comprehensive server tests
- `apps/core/cmd/autobmad/main.go` - Wired up server with graceful shutdown

## Dev Agent Record

### Agent Model Used

Claude 3.5 Sonnet (claude-sonnet-4-20250514)

### Completion Notes List

- Implemented JSON-RPC 2.0 compliant server with length-prefixed framing
- Custom UnmarshalJSON for Request to distinguish between `"id": null` and no id field (important for notification detection per spec)
- All 5 standard error codes implemented and tested: -32700, -32600, -32601, -32602, -32603
- System handlers registered: system.ping (returns "pong"), system.echo (echoes message)
- Server supports graceful shutdown via context cancellation and SIGTERM/SIGINT
- All logging goes to stderr; stdout is reserved for JSON-RPC messages only
- 55 tests total, all passing
- Frame format: [4-byte big-endian length][JSON payload][newline]

### Change Log

| Date | Change | Reason |
|------|--------|--------|
| 2026-01-22 | Initial implementation | Story 1.2 implementation
