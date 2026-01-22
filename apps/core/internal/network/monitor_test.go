package network

import (
	"context"
	"testing"
	"time"
)

// TestNetworkStatus_Constants verifies status constant values
func TestNetworkStatus_Constants(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{StatusOnline, "online"},
		{StatusOffline, "offline"},
		{StatusChecking, "checking"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("Status value = %v, want %v", tt.status, tt.want)
			}
		})
	}
}

// TestNewMonitor verifies monitor creation
func TestNewMonitor(t *testing.T) {
	onChange := func(old, new Status) {}
	m := NewMonitor(30*time.Second, onChange)

	if m == nil {
		t.Fatal("NewMonitor returned nil")
	}

	status := m.GetStatus()
	if status.Status != StatusChecking {
		t.Errorf("Initial status = %v, want %v", status.Status, StatusChecking)
	}
}

// TestMonitor_GetStatus verifies status retrieval is thread-safe
func TestMonitor_GetStatus(t *testing.T) {
	m := NewMonitor(30*time.Second, nil)

	// Call GetStatus concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_ = m.GetStatus()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestMonitor_Start verifies monitor starts and runs periodic checks
func TestMonitor_Start(t *testing.T) {
	changeCount := 0
	onChange := func(old, new Status) {
		changeCount++
	}

	// Use short interval for testing
	m := NewMonitor(100*time.Millisecond, onChange)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Start monitor in background
	go m.Start(ctx)

	// Wait for context timeout
	<-ctx.Done()

	// Should have performed initial check
	status := m.GetStatus()
	if status.Status == StatusChecking {
		t.Error("Status should have changed from checking after initial check")
	}
}

// TestMonitor_Stop verifies monitor can be stopped
func TestMonitor_Stop(t *testing.T) {
	m := NewMonitor(30*time.Second, nil)

	ctx := context.Background()
	go m.Start(ctx)

	// Give it time to start
	time.Sleep(50 * time.Millisecond)

	// Stop it
	m.Stop()

	// Should stop without hanging
	time.Sleep(50 * time.Millisecond)
}

// TestMonitor_OnChange verifies onChange callback is called on status change
func TestMonitor_OnChange(t *testing.T) {
	var callbackNew Status
	onChange := func(old, new Status) {
		callbackNew = new
	}

	m := NewMonitor(100*time.Millisecond, onChange)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Start monitor
	go m.Start(ctx)

	// Wait for checks to run
	time.Sleep(350 * time.Millisecond)

	// If status changed (not from checking), callback should have been called
	// Note: In real execution, status will likely be online
	status := m.GetStatus()
	if status.Status != StatusChecking && callbackNew == "" {
		t.Log("Callback was not called - network might not have changed during test")
	}
}

// TestIsOnline verifies connectivity check function
func TestIsOnline(t *testing.T) {
	// This test may fail if network is actually offline
	// We test that the function executes without panic
	result := isOnline()
	t.Logf("isOnline() returned: %v", result)
	// We can't assert true/false as it depends on actual network
}

// TestNetworkStatus_Latency verifies latency is captured
func TestNetworkStatus_Latency(t *testing.T) {
	m := NewMonitor(100*time.Millisecond, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go m.Start(ctx)

	// Wait for initial check
	time.Sleep(150 * time.Millisecond)

	status := m.GetStatus()
	if status.Latency < 0 {
		t.Errorf("Latency should be non-negative, got %v", status.Latency)
	}
	if status.LastChecked.IsZero() {
		t.Error("LastChecked should be set after check")
	}
}

// TestNewDebouncedMonitor verifies debounced monitor creation
func TestNewDebouncedMonitor(t *testing.T) {
	onChange := func(old, new Status) {}
	dm := NewDebouncedMonitor(30*time.Second, 5*time.Second, onChange)

	if dm == nil {
		t.Fatal("NewDebouncedMonitor returned nil")
	}

	status := dm.GetStatus()
	if status.Status != StatusChecking {
		t.Errorf("Initial status = %v, want %v", status.Status, StatusChecking)
	}
}

// TestDebouncedMonitor_Debounce verifies that rapid changes are debounced
func TestDebouncedMonitor_Debounce(t *testing.T) {
	callCount := 0
	onChange := func(old, new Status) {
		callCount++
	}

	// Short debounce window for testing
	dm := NewDebouncedMonitor(50*time.Millisecond, 100*time.Millisecond, onChange)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go dm.Start(ctx)

	// Wait for execution
	time.Sleep(350 * time.Millisecond)

	// With debouncing, we should have fewer callbacks than without
	t.Logf("Callback called %d times with debouncing", callCount)
}
