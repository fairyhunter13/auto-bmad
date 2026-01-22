// Package server provides the JSON-RPC 2.0 server over stdio.
// This package handles incoming requests from the Electron main process
// and routes them to the appropriate handlers.
package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"sync"
)

// Handler is a function that processes JSON-RPC method calls.
// It receives the params as raw JSON and returns a result or error.
type Handler func(params json.RawMessage) (interface{}, error)

// Server represents the JSON-RPC server that communicates over stdio.
type Server struct {
	reader   *MessageReader
	writer   *MessageWriter
	handlers map[string]Handler
	logger   *log.Logger
	mu       sync.RWMutex // protects handlers map
}

// New creates a new JSON-RPC server instance.
// stdin and stdout are the I/O streams for JSON-RPC communication.
// logger should write to stderr (stdout is reserved for JSON-RPC).
func New(stdin io.Reader, stdout io.Writer, logger *log.Logger) *Server {
	return &Server{
		reader:   NewMessageReader(stdin),
		writer:   NewMessageWriter(stdout),
		handlers: make(map[string]Handler),
		logger:   logger,
	}
}

// RegisterHandler registers a handler for the given method name.
// Handler names should follow the "resource.action" convention.
func (s *Server) RegisterHandler(method string, handler Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[method] = handler
}

// EmitEvent sends a server-initiated notification to the client.
// Event names should follow the "resource.event" convention.
// This sends a JSON-RPC notification (no ID, no response expected).
func (s *Server) EmitEvent(event string, data interface{}) error {
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  event,
		"params":  data,
	}

	// Write using the framing protocol
	if err := s.writer.writeJSON(notification); err != nil {
		return err
	}

	s.logger.Printf("Event emitted: %s", event)
	return nil
}

// Run starts the server loop, processing requests until the context is cancelled
// or stdin is closed (EOF). Returns nil on clean shutdown, or an error otherwise.
func (s *Server) Run(ctx context.Context) error {
	// Create a channel for read results
	type readResult struct {
		req *Request
		err error
	}
	readCh := make(chan readResult, 1)

	for {
		// Start a read in a goroutine so we can also check context cancellation
		go func() {
			req, err := s.reader.ReadRequest()
			readCh <- readResult{req, err}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()

		case result := <-readCh:
			if result.err != nil {
				if result.err == io.EOF {
					return nil // Clean shutdown - stdin closed
				}
				// Parse error - invalid JSON framing or JSON syntax
				s.writeParseError(result.err)
				continue
			}

			s.handleRequest(result.req)
		}
	}
}

// handleRequest processes a single JSON-RPC request.
func (s *Server) handleRequest(req *Request) {
	s.logger.Printf("Request: method=%s id=%v", req.Method, req.ID)

	// Validate JSON-RPC 2.0 request
	if !req.IsValid() {
		if !req.IsNotification() {
			s.writeError(req.ID, ErrCodeInvalidRequest, "Invalid Request", "jsonrpc must be \"2.0\" and method must be non-empty")
		}
		return
	}

	// Find handler
	s.mu.RLock()
	handler, ok := s.handlers[req.Method]
	s.mu.RUnlock()

	if !ok {
		if !req.IsNotification() {
			s.writeError(req.ID, ErrCodeMethodNotFound, "Method not found", req.Method)
		}
		return
	}

	// Execute handler
	result, err := handler(req.Params)
	if err != nil {
		if !req.IsNotification() {
			// Check if it's already a JSON-RPC Error
			if rpcErr, ok := err.(*Error); ok {
				s.writeErrorResponse(req.ID, rpcErr)
			} else {
				s.writeError(req.ID, ErrCodeInternalError, "Internal error", err.Error())
			}
		}
		return
	}

	// Write success response (only for requests, not notifications)
	if !req.IsNotification() {
		s.writeResult(req.ID, result)
	}
}

// writeResult writes a success response.
func (s *Server) writeResult(id interface{}, result interface{}) {
	resp := NewSuccessResponse(id, result)
	if err := s.writer.WriteResponse(resp); err != nil {
		s.logger.Printf("Error writing response: %v", err)
	}
	s.logger.Printf("Response: id=%v result=%v", id, result)
}

// writeError writes an error response with the given code and message.
func (s *Server) writeError(id interface{}, code int, message string, data interface{}) {
	resp := NewErrorResponseWithData(id, code, message, data)
	if err := s.writer.WriteResponse(resp); err != nil {
		s.logger.Printf("Error writing error response: %v", err)
	}
	s.logger.Printf("Response: id=%v error=%d %s", id, code, message)
}

// writeErrorResponse writes an error response from an Error struct.
func (s *Server) writeErrorResponse(id interface{}, rpcErr *Error) {
	resp := &Response{
		JSONRPC: "2.0",
		Error:   rpcErr,
		ID:      id,
	}
	if err := s.writer.WriteResponse(resp); err != nil {
		s.logger.Printf("Error writing error response: %v", err)
	}
	s.logger.Printf("Response: id=%v error=%d %s", id, rpcErr.Code, rpcErr.Message)
}

// writeParseError writes a parse error response (used when JSON parsing fails).
// Per JSON-RPC 2.0 spec, parse errors have null ID since we couldn't parse the request.
func (s *Server) writeParseError(parseErr error) {
	s.logger.Printf("Parse error: %v", parseErr)
	resp := NewErrorResponseWithData(nil, ErrCodeParseError, "Parse error", parseErr.Error())
	if err := s.writer.WriteResponse(resp); err != nil {
		s.logger.Printf("Error writing parse error response: %v", err)
	}
}
