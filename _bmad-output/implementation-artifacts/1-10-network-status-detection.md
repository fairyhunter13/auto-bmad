# Story 1.10: Network Status Detection

Status: ready-for-dev

## Story

As a **user running journeys with cloud AI providers**,
I want **Auto-BMAD to monitor network connectivity and warn me of issues**,
So that **I understand why journeys might fail and can take action**.

## Acceptance Criteria

1. **Given** the user is online  
   **When** Auto-BMAD checks network status  
   **Then** status indicator shows "connected"  
   **And** cloud provider workflows are enabled

2. **Given** the user goes offline  
   **When** network status changes  
   **Then** a non-intrusive warning notification appears  
   **And** status indicator updates to "offline"  
   **And** if a journey is running, user is asked: "Pause journey or continue anyway?"

3. **Given** the user is using a local provider (e.g., Ollama)  
   **When** network status is offline  
   **Then** journeys can still execute  
   **And** UI indicates "Using local provider - offline operation available"

## Tasks / Subtasks

- [ ] **Task 1: Implement network status check in Golang** (AC: #1)
  - [ ] Create `internal/network/monitor.go`
  - [ ] Implement connectivity check (DNS lookup or HTTP ping)
  - [ ] Return structured status result
  - [ ] Run check periodically (every 30 seconds)

- [ ] **Task 2: Implement real-time status events** (AC: #2)
  - [ ] Emit `network.statusChanged` event on state change
  - [ ] Include previous and new status in event
  - [ ] Debounce rapid changes (5 second window)

- [ ] **Task 3: Create JSON-RPC handler and event** (AC: #1, #2)
  - [ ] Register `network.getStatus` method
  - [ ] Set up periodic status broadcast
  - [ ] Handle pause prompt for active journeys

- [ ] **Task 4: Add frontend status indicator** (AC: all)
  - [ ] Create NetworkStatusIndicator component
  - [ ] Add to status bar or header
  - [ ] Show warning modal on offline transition
  - [ ] Handle local provider detection

## Dev Notes

### Network Monitoring Implementation

```go
// internal/network/monitor.go

package network

import (
    "context"
    "net"
    "sync"
    "time"
)

type Status string

const (
    StatusOnline   Status = "online"
    StatusOffline  Status = "offline"
    StatusChecking Status = "checking"
)

type NetworkStatus struct {
    Status      Status    `json:"status"`
    LastChecked time.Time `json:"lastChecked"`
    Latency     int64     `json:"latency,omitempty"` // ms
}

type Monitor struct {
    status     NetworkStatus
    mu         sync.RWMutex
    interval   time.Duration
    onChange   func(old, new Status)
    stopCh     chan struct{}
}

func NewMonitor(interval time.Duration, onChange func(old, new Status)) *Monitor {
    return &Monitor{
        status:   NetworkStatus{Status: StatusChecking},
        interval: interval,
        onChange: onChange,
        stopCh:   make(chan struct{}),
    }
}

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

func (m *Monitor) Stop() {
    close(m.stopCh)
}

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
    
    // Notify on change
    if oldStatus != newStatus && oldStatus != StatusChecking {
        if m.onChange != nil {
            m.onChange(oldStatus, newStatus)
        }
    }
}

func (m *Monitor) GetStatus() NetworkStatus {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.status
}

func isOnline() bool {
    // Method 1: DNS lookup (fast, reliable)
    _, err := net.LookupHost("dns.google")
    if err == nil {
        return true
    }
    
    // Method 2: Fallback to another DNS
    _, err = net.LookupHost("cloudflare.com")
    return err == nil
}
```

### Event Emission with Debounce

```go
// internal/network/monitor.go (additions)

type DebouncedMonitor struct {
    *Monitor
    debounceWindow time.Duration
    lastChange     time.Time
    pendingNotify  *time.Timer
    mu             sync.Mutex
}

func NewDebouncedMonitor(interval, debounce time.Duration, onChange func(old, new Status)) *DebouncedMonitor {
    dm := &DebouncedMonitor{
        debounceWindow: debounce,
    }
    
    dm.Monitor = NewMonitor(interval, func(old, new Status) {
        dm.debouncedNotify(old, new, onChange)
    })
    
    return dm
}

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
```

### JSON-RPC Handler and Events

```go
// internal/server/handlers.go (additions)

func (s *Server) handleNetworkGetStatus(params json.RawMessage) (interface{}, error) {
    return s.networkMonitor.GetStatus(), nil
}

// In server initialization
s.networkMonitor = network.NewDebouncedMonitor(
    30*time.Second,  // Check interval
    5*time.Second,   // Debounce window
    func(old, new network.Status) {
        // Emit event to frontend
        s.emitEvent("network.statusChanged", map[string]interface{}{
            "previous": old,
            "current":  new,
        })
    },
)

// Start monitoring
go s.networkMonitor.Start(ctx)
```

### Frontend Types

```typescript
// src/renderer/types/network.ts

export type NetworkStatus = 'online' | 'offline' | 'checking';

export interface NetworkStatusResult {
  status: NetworkStatus;
  lastChecked: string;
  latency?: number;
}

export interface NetworkStatusChangedEvent {
  previous: NetworkStatus;
  current: NetworkStatus;
}
```

### Preload API Additions

```typescript
// src/preload/index.ts (additions)

const api = {
  network: {
    getStatus: (): Promise<NetworkStatusResult> => 
      ipcRenderer.invoke('rpc:call', 'network.getStatus'),
    onStatusChanged: (callback: (event: NetworkStatusChangedEvent) => void) => {
      ipcRenderer.on('event:network.statusChanged', (_, event) => callback(event));
    },
  },
};
```

### Network Status Indicator Component

```tsx
// src/renderer/components/NetworkStatusIndicator.tsx

import { useEffect, useState } from 'react';
import { Wifi, WifiOff } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';

export function NetworkStatusIndicator() {
  const [status, setStatus] = useState<NetworkStatus>('checking');
  const [showOfflineDialog, setShowOfflineDialog] = useState(false);
  const [hasActiveJourney, setHasActiveJourney] = useState(false);
  
  useEffect(() => {
    // Get initial status
    window.api.network.getStatus().then((result) => {
      setStatus(result.status);
    });
    
    // Listen for changes
    window.api.network.onStatusChanged((event) => {
      setStatus(event.current);
      
      if (event.current === 'offline' && event.previous === 'online') {
        // Check if journey is active
        // If so, show dialog
        if (hasActiveJourney) {
          setShowOfflineDialog(true);
        }
      }
    });
  }, [hasActiveJourney]);
  
  return (
    <>
      <Badge 
        variant={status === 'online' ? 'default' : 'destructive'}
        className="flex items-center gap-1"
      >
        {status === 'online' ? (
          <>
            <Wifi className="h-3 w-3" />
            <span>Online</span>
          </>
        ) : (
          <>
            <WifiOff className="h-3 w-3" />
            <span>Offline</span>
          </>
        )}
      </Badge>
      
      {/* Offline Warning Dialog */}
      <AlertDialog open={showOfflineDialog} onOpenChange={setShowOfflineDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Network Connection Lost</AlertDialogTitle>
            <AlertDialogDescription>
              You appear to be offline. Cloud AI providers may not be accessible.
              Would you like to pause the current journey or continue anyway?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => handlePauseJourney()}>
              Pause Journey
            </AlertDialogCancel>
            <AlertDialogAction onClick={() => setShowOfflineDialog(false)}>
              Continue Anyway
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
```

### Local Provider Detection

```typescript
// src/renderer/hooks/useProviderType.ts

export function useProviderType(selectedProfile: string) {
  // Detect if profile uses local provider
  const isLocalProvider = useMemo(() => {
    const localPatterns = ['ollama', 'local', 'localhost'];
    return localPatterns.some(pattern => 
      selectedProfile.toLowerCase().includes(pattern)
    );
  }, [selectedProfile]);
  
  return { isLocalProvider };
}

// In NetworkStatusIndicator
const { isLocalProvider } = useProviderType(selectedProfile);

// If offline but using local provider
if (status === 'offline' && isLocalProvider) {
  return (
    <Badge variant="secondary" className="flex items-center gap-1">
      <Wifi className="h-3 w-3" />
      <span>Local Mode</span>
    </Badge>
  );
}
```

### File Structure

```
apps/core/internal/
└── network/
    ├── monitor.go       # Network monitoring
    └── monitor_test.go  # Unit tests

apps/desktop/src/renderer/
├── components/
│   └── NetworkStatusIndicator.tsx
└── hooks/
    └── useProviderType.ts
```

### Status Bar Integration

Add NetworkStatusIndicator to the main layout status bar:

```tsx
// src/renderer/layouts/MainLayout.tsx

export function MainLayout({ children }) {
  return (
    <div className="flex flex-col h-screen">
      {/* Main content */}
      <main className="flex-1">{children}</main>
      
      {/* Status bar */}
      <footer className="border-t px-4 py-2 flex items-center justify-between text-sm">
        <NetworkStatusIndicator />
        <span className="text-muted-foreground">Auto-BMAD v1.0</span>
      </footer>
    </div>
  );
}
```

### Testing Requirements

1. Test online detection (DNS lookup success)
2. Test offline detection (DNS lookup failure)
3. Test debounced event emission
4. Test status change notification
5. Test local provider bypass

### Dependencies

- **Story 1.3**: IPC bridge for events
- **Story 1.5**: Profile detection for local provider check

### References

- [prd.md#FR51] - System can detect network connectivity status
- [prd.md#FR52] - System can warn user of network issues during journey

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Completion Notes List

- 

### Change Log

| Date | Change | Reason |
|------|--------|--------|
