// Package main provides integration tests for the autobmad binary.
package main

import (
	"encoding/binary"
	"encoding/json"
	"os/exec"
	"testing"
	"time"

	"github.com/fairyhunter13/auto-bmad/apps/core/internal/server"
)

func TestBinaryPingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", "autobmad_test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer exec.Command("rm", "autobmad_test").Run()

	// Start the binary with test project path
	testProjectPath := t.TempDir()
	proc := exec.Command("./autobmad_test", "--project-path", testProjectPath)
	stdin, err := proc.StdinPipe()
	if err != nil {
		t.Fatalf("failed to get stdin pipe: %v", err)
	}
	stdout, err := proc.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to get stdout pipe: %v", err)
	}

	if err := proc.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}

	// Send a ping request
	request := `{"jsonrpc":"2.0","method":"system.ping","id":1}`
	frame := make([]byte, 4+len(request)+1)
	binary.BigEndian.PutUint32(frame[:4], uint32(len(request)))
	copy(frame[4:], request)
	frame[len(frame)-1] = '\n'

	if _, err := stdin.Write(frame); err != nil {
		t.Fatalf("failed to write request: %v", err)
	}

	// Read response with timeout
	done := make(chan struct{})
	var resp server.Response
	var readErr error

	go func() {
		defer close(done)
		lengthBuf := make([]byte, 4)
		if _, err := stdout.Read(lengthBuf); err != nil {
			readErr = err
			return
		}
		length := binary.BigEndian.Uint32(lengthBuf)
		payload := make([]byte, length)
		if _, err := stdout.Read(payload); err != nil {
			readErr = err
			return
		}
		readErr = json.Unmarshal(payload, &resp)
	}()

	select {
	case <-done:
		if readErr != nil {
			t.Fatalf("failed to read response: %v", readErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for response")
	}

	// Close stdin to signal EOF and wait for process
	stdin.Close()
	proc.Wait()

	// Verify response
	if resp.Result != "pong" {
		t.Errorf("expected 'pong', got %v", resp.Result)
	}
}

func TestBinaryGracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", "autobmad_test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer exec.Command("rm", "autobmad_test").Run()

	// Start the binary with test project path
	testProjectPath := t.TempDir()
	proc := exec.Command("./autobmad_test", "--project-path", testProjectPath)
	stdin, _ := proc.StdinPipe()

	if err := proc.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Close stdin to trigger EOF shutdown
	stdin.Close()

	// Wait for process with timeout
	done := make(chan error, 1)
	go func() {
		done <- proc.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("process exited with error: %v", err)
		}
	case <-time.After(5 * time.Second):
		proc.Process.Kill()
		t.Fatal("process did not shutdown within timeout")
	}
}
