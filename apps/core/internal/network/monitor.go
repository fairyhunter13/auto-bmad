// Package network provides network connectivity monitoring.
package network

import (
	"context"
	"net"
	"sync"
	"time"
)

// Status represents the current network connectivity status.
type Status string

const (
	// StatusOnline indicates network is available
	StatusOnline Status = "online"
	// StatusOffline indicates network is unavailable
	StatusOffline Status = "offline"
	// StatusChecking indicates status check is in progress
	StatusChecking Status = "checking"
)

// NetworkStatus represents the result of a network connectivity check.
type NetworkStatus struct {
	Status      Status    `json:"status"`
	LastChecked time.Time `json:"lastChecked"`
	Latency     int64     `json:"latency,omitempty"` // milliseconds
}

// Monitor continuously monitors network connectivity status.
type Monitor struct {
	status   NetworkStatus
	mu       sync.RWMutex
	interval time.Duration
	onChange func(old, new Status)
	stopCh   chan struct{}
}

// NewMonitor creates a new network monitor that checks connectivity at the given interval.
// The onChange callback is called whenever the status changes (not including initial check).
func NewMonitor(interval time.Duration, onChange func(old, new Status)) *Monitor {
	return &Monitor{
		status:   NetworkStatus{Status: StatusChecking},
		interval: interval,
		onChange: onChange,
		stopCh:   make(chan struct{}),
	}
}

// Start begins monitoring network status until the context is cancelled or Stop is called.
// It performs an initial check immediately, then continues at the configured interval.
func (m *Monitor) Start(ctx context.Context) {
	// Initial check
	m.check()

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.check()
		case <-ctx.Done():
			return
		case <-m.stopCh:
			return
		}
	}
}

// Stop stops the network monitor.
func (m *Monitor) Stop() {
	close(m.stopCh)
}

// check performs a single network connectivity check and updates status.
func (m *Monitor) check() {
	m.mu.Lock()
	oldStatus := m.status.Status
	m.mu.Unlock()

	start := time.Now()
	online := isOnline()
	latency := time.Since(start).Milliseconds()

	newStatus := StatusOffline
	if online {
		newStatus = StatusOnline
	}

	m.mu.Lock()
	m.status = NetworkStatus{
		Status:      newStatus,
		LastChecked: time.Now(),
		Latency:     latency,
	}
	m.mu.Unlock()

	// Notify on change (skip initial check from StatusChecking)
	if oldStatus != newStatus && oldStatus != StatusChecking {
		if m.onChange != nil {
			m.onChange(oldStatus, newStatus)
		}
	}
}

// GetStatus returns the current network status.
// This method is thread-safe.
func (m *Monitor) GetStatus() NetworkStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status
}

// isOnline checks if the system has network connectivity.
// It uses DNS lookups to well-known hosts as a connectivity test.
func isOnline() bool {
	// Method 1: DNS lookup to Google DNS (fast, reliable)
	_, err := net.LookupHost("dns.google")
	if err == nil {
		return true
	}

	// Method 2: Fallback to another DNS provider
	_, err = net.LookupHost("cloudflare.com")
	return err == nil
}

// DebouncedMonitor wraps Monitor with debouncing for status change notifications.
// This prevents rapid status changes from triggering too many callbacks.
type DebouncedMonitor struct {
	*Monitor
	debounceWindow time.Duration
	lastChange     time.Time
	pendingNotify  *time.Timer
	mu             sync.Mutex
}

// NewDebouncedMonitor creates a new debounced network monitor.
// Status changes are debounced by the given window duration.
func NewDebouncedMonitor(interval, debounce time.Duration, onChange func(old, new Status)) *DebouncedMonitor {
	dm := &DebouncedMonitor{
		debounceWindow: debounce,
	}

	dm.Monitor = NewMonitor(interval, func(old, new Status) {
		dm.debouncedNotify(old, new, onChange)
	})

	return dm
}

// debouncedNotify delays the onChange callback by the debounce window.
// If another change occurs within the window, the timer is reset.
func (dm *DebouncedMonitor) debouncedNotify(old, new Status, callback func(Status, Status)) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Cancel pending notification
	if dm.pendingNotify != nil {
		dm.pendingNotify.Stop()
	}

	// Schedule debounced notification
	dm.pendingNotify = time.AfterFunc(dm.debounceWindow, func() {
		callback(old, new)
	})
}
