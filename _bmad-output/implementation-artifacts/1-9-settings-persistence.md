# Story 1.9: Settings Persistence

Status: ready-for-dev

## Story

As a **user who has configured Auto-BMAD**,
I want **my settings to be saved and restored across sessions**,
So that **I don't have to reconfigure every time I open the app**.

## Acceptance Criteria

1. **Given** the user changes settings (retry limits, notification preferences)  
   **When** the changes are made  
   **Then** settings are saved to `_bmad-output/.autobmad/config.json`  
   **And** the save completes within 1 second (NFR-P6)

2. **Given** the user restarts Auto-BMAD  
   **When** the app initializes  
   **Then** previous settings are restored automatically  
   **And** the last-used project folder is remembered  
   **And** the last-used OpenCode profile per project is restored

3. **Given** no settings file exists  
   **When** the app initializes  
   **Then** sensible defaults are applied  
   **And** a new settings file is created on first change

## Tasks / Subtasks

- [ ] **Task 1: Define settings schema** (AC: #1, #3)
  - [ ] Create `internal/state/config.go` with settings struct
  - [ ] Define default values for all settings
  - [ ] Include retry limits, notification prefs, recent projects

- [ ] **Task 2: Implement settings persistence** (AC: #1, #2)
  - [ ] Create `internal/state/manager.go` for state management
  - [ ] Implement save to JSON with atomic write
  - [ ] Implement load from JSON with defaults fallback
  - [ ] Ensure save completes within 1 second

- [ ] **Task 3: Create JSON-RPC handlers** (AC: all)
  - [ ] Register `settings.get` method
  - [ ] Register `settings.set` method
  - [ ] Register `settings.reset` method

- [ ] **Task 4: Add frontend API and settings UI** (AC: all)
  - [ ] Add settings API to preload
  - [ ] Create Settings screen component
  - [ ] Implement form with shadcn/ui components

## Dev Notes

### Architecture Requirements

**Source:** [architecture.md#State Architecture]

| Aspect | Decision |
|--------|----------|
| Configuration | `_bmad-output/.autobmad/config.json` |
| Crash Recovery | Read filesystem state on restart |
| Save Time | < 1 second (NFR-P6) |

### Settings Schema

```go
// internal/state/config.go

package state

import (
    "encoding/json"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Settings struct {
    // Retry settings
    MaxRetries   int `json:"maxRetries"`   // Default: 3
    RetryDelay   int `json:"retryDelay"`   // Default: 5000 (ms)
    
    // Notification settings
    DesktopNotifications bool `json:"desktopNotifications"` // Default: true
    SoundEnabled         bool `json:"soundEnabled"`         // Default: false
    
    // Timeout settings
    StepTimeoutDefault   int `json:"stepTimeoutDefault"`   // Default: 300000 (5 min)
    HeartbeatInterval    int `json:"heartbeatInterval"`    // Default: 60000 (60s)
    
    // UI preferences
    Theme           string `json:"theme"`           // Default: "system"
    ShowDebugOutput bool   `json:"showDebugOutput"` // Default: false
    
    // Project memory
    LastProjectPath      string            `json:"lastProjectPath,omitempty"`
    ProjectProfiles      map[string]string `json:"projectProfiles"`      // path -> profile name
    RecentProjectsMax    int               `json:"recentProjectsMax"`    // Default: 10
}

func DefaultSettings() *Settings {
    return &Settings{
        MaxRetries:           3,
        RetryDelay:           5000,
        DesktopNotifications: true,
        SoundEnabled:         false,
        StepTimeoutDefault:   300000,
        HeartbeatInterval:    60000,
        Theme:                "system",
        ShowDebugOutput:      false,
        ProjectProfiles:      make(map[string]string),
        RecentProjectsMax:    10,
    }
}
```

### State Manager Implementation

```go
// internal/state/manager.go

type StateManager struct {
    settings    *Settings
    configPath  string
    mu          sync.RWMutex
    dirty       bool
    saveTimeout *time.Timer
}

func NewStateManager(projectPath string) (*StateManager, error) {
    configDir := filepath.Join(projectPath, "_bmad-output", ".autobmad")
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return nil, err
    }
    
    sm := &StateManager{
        configPath: filepath.Join(configDir, "config.json"),
        settings:   DefaultSettings(),
    }
    
    // Load existing settings
    if err := sm.load(); err != nil && !os.IsNotExist(err) {
        return nil, err
    }
    
    return sm, nil
}

func (sm *StateManager) load() error {
    data, err := os.ReadFile(sm.configPath)
    if err != nil {
        return err
    }
    
    settings := DefaultSettings() // Start with defaults
    if err := json.Unmarshal(data, settings); err != nil {
        return err
    }
    
    sm.settings = settings
    return nil
}

func (sm *StateManager) save() error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    data, err := json.MarshalIndent(sm.settings, "", "  ")
    if err != nil {
        return err
    }
    
    // Atomic write: write to temp file, then rename
    tempPath := sm.configPath + ".tmp"
    if err := os.WriteFile(tempPath, data, 0644); err != nil {
        return err
    }
    
    return os.Rename(tempPath, sm.configPath)
}

func (sm *StateManager) Get() *Settings {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    // Return a copy to prevent mutation
    copy := *sm.settings
    return &copy
}

func (sm *StateManager) Set(updates map[string]interface{}) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    // Apply updates
    for key, value := range updates {
        switch key {
        case "maxRetries":
            sm.settings.MaxRetries = int(value.(float64))
        case "retryDelay":
            sm.settings.RetryDelay = int(value.(float64))
        case "desktopNotifications":
            sm.settings.DesktopNotifications = value.(bool)
        case "soundEnabled":
            sm.settings.SoundEnabled = value.(bool)
        case "theme":
            sm.settings.Theme = value.(string)
        case "showDebugOutput":
            sm.settings.ShowDebugOutput = value.(bool)
        // ... other fields
        }
    }
    
    sm.dirty = true
    return sm.save()
}

func (sm *StateManager) Reset() error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    sm.settings = DefaultSettings()
    return sm.save()
}
```

### JSON-RPC Handlers

```go
// internal/server/handlers.go (additions)

func (s *Server) handleSettingsGet(params json.RawMessage) (interface{}, error) {
    return s.stateManager.Get(), nil
}

func (s *Server) handleSettingsSet(params json.RawMessage) (interface{}, error) {
    var updates map[string]interface{}
    if err := json.Unmarshal(params, &updates); err != nil {
        return nil, err
    }
    
    if err := s.stateManager.Set(updates); err != nil {
        return nil, err
    }
    
    return s.stateManager.Get(), nil
}

func (s *Server) handleSettingsReset(params json.RawMessage) (interface{}, error) {
    if err := s.stateManager.Reset(); err != nil {
        return nil, err
    }
    return s.stateManager.Get(), nil
}

// Register handlers
server.RegisterHandler("settings.get", s.handleSettingsGet)
server.RegisterHandler("settings.set", s.handleSettingsSet)
server.RegisterHandler("settings.reset", s.handleSettingsReset)
```

### Frontend Types

```typescript
// src/renderer/types/settings.ts

export interface Settings {
  maxRetries: number;
  retryDelay: number;
  desktopNotifications: boolean;
  soundEnabled: boolean;
  stepTimeoutDefault: number;
  heartbeatInterval: number;
  theme: 'light' | 'dark' | 'system';
  showDebugOutput: boolean;
  lastProjectPath?: string;
  projectProfiles: Record<string, string>;
  recentProjectsMax: number;
}
```

### Preload API Additions

```typescript
// src/preload/index.ts (additions)

const api = {
  settings: {
    get: (): Promise<Settings> => 
      ipcRenderer.invoke('rpc:call', 'settings.get'),
    set: (updates: Partial<Settings>): Promise<Settings> => 
      ipcRenderer.invoke('rpc:call', 'settings.set', updates),
    reset: (): Promise<Settings> => 
      ipcRenderer.invoke('rpc:call', 'settings.reset'),
  },
};
```

### Settings Screen Component

```tsx
// src/renderer/screens/SettingsScreen.tsx

import { useEffect, useState } from 'react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export function SettingsScreen() {
  const [settings, setSettings] = useState<Settings | null>(null);
  const [saving, setSaving] = useState(false);
  
  useEffect(() => {
    window.api.settings.get().then(setSettings);
  }, []);
  
  const handleChange = async (key: keyof Settings, value: unknown) => {
    setSaving(true);
    try {
      const updated = await window.api.settings.set({ [key]: value });
      setSettings(updated);
    } finally {
      setSaving(false);
    }
  };
  
  if (!settings) return <div>Loading...</div>;
  
  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8">Settings</h1>
      
      <div className="space-y-6">
        {/* Retry Settings */}
        <Card>
          <CardHeader>
            <CardTitle>Retry Behavior</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="maxRetries">Maximum Retries</Label>
              <Input
                id="maxRetries"
                type="number"
                value={settings.maxRetries}
                onChange={(e) => handleChange('maxRetries', parseInt(e.target.value))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="retryDelay">Retry Delay (ms)</Label>
              <Input
                id="retryDelay"
                type="number"
                value={settings.retryDelay}
                onChange={(e) => handleChange('retryDelay', parseInt(e.target.value))}
              />
            </div>
          </CardContent>
        </Card>
        
        {/* Notification Settings */}
        <Card>
          <CardHeader>
            <CardTitle>Notifications</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <Label htmlFor="desktopNotifications">Desktop Notifications</Label>
              <Switch
                id="desktopNotifications"
                checked={settings.desktopNotifications}
                onCheckedChange={(v) => handleChange('desktopNotifications', v)}
              />
            </div>
            <div className="flex items-center justify-between">
              <Label htmlFor="soundEnabled">Sound Effects</Label>
              <Switch
                id="soundEnabled"
                checked={settings.soundEnabled}
                onCheckedChange={(v) => handleChange('soundEnabled', v)}
              />
            </div>
          </CardContent>
        </Card>
        
        {/* Reset Button */}
        <Button 
          variant="destructive" 
          onClick={async () => {
            if (confirm('Reset all settings to defaults?')) {
              const defaults = await window.api.settings.reset();
              setSettings(defaults);
            }
          }}
        >
          Reset to Defaults
        </Button>
      </div>
    </div>
  );
}
```

### File Structure

```
apps/core/internal/state/
├── config.go         # Settings struct and defaults
├── manager.go        # StateManager implementation
└── manager_test.go   # Unit tests

apps/desktop/src/renderer/screens/
└── SettingsScreen.tsx
```

### Performance Requirement

**Source:** [prd.md#NFR-P6]

Settings save MUST complete within 1 second. Atomic write ensures no data loss.

### Testing Requirements

1. Test default settings applied on first run
2. Test settings persistence across restarts
3. Test atomic write (no corruption on crash)
4. Test save time < 1 second
5. Test project-specific profile memory

### Dependencies

- **Story 1.7**: Project detection for config path
- **Story 1.3**: IPC bridge must be working

### References

- [architecture.md#State Architecture] - Configuration location
- [prd.md#FR49] - User can configure Auto-BMAD settings
- [prd.md#FR50] - System can persist user preferences
- [prd.md#NFR-P6] - Journey state save < 1 second

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Completion Notes List

- 

### Change Log

| Date | Change | Reason |
|------|--------|--------|
