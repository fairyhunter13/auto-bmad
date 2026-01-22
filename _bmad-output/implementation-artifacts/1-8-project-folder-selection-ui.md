# Story 1.8: Project Folder Selection UI

Status: ready-for-dev

## Story

As a **user starting Auto-BMAD**,
I want **to select a project folder and optionally describe its context**,
So that **Auto-BMAD knows which project to work with**.

## Acceptance Criteria

1. **Given** the user opens Auto-BMAD  
   **When** no project is selected  
   **Then** a project selection screen is displayed  
   **And** a "Select Project Folder" button opens the native file picker  
   **And** recently opened projects are listed for quick access

2. **Given** the user selects a valid BMAD project folder  
   **When** the selection is confirmed  
   **Then** project detection runs automatically  
   **And** results are displayed (BMAD version, greenfield/brownfield, dependencies)  
   **And** the project is saved to recent projects list

3. **Given** a brownfield project is selected  
   **When** the user wants to provide additional context  
   **Then** a text area allows manual project description (FR48)  
   **And** the description is saved with the project configuration

## Tasks / Subtasks

- [ ] **Task 1: Create project selection screen** (AC: #1)
  - [ ] Create `ProjectSelectScreen` component
  - [ ] Add "Select Project Folder" button with native dialog
  - [ ] Display recently opened projects list
  - [ ] Style with shadcn/ui Card and Button components

- [ ] **Task 2: Implement native folder picker** (AC: #1, #2)
  - [ ] Add `window.api.dialog.selectFolder()` to preload
  - [ ] Use Electron's `dialog.showOpenDialog` with directory mode
  - [ ] Return selected path to renderer

- [ ] **Task 3: Display project detection results** (AC: #2)
  - [ ] Create `ProjectInfo` component
  - [ ] Show BMAD version and compatibility
  - [ ] Show project type (greenfield/brownfield)
  - [ ] Show dependency status (OpenCode, Git)
  - [ ] List existing artifacts for brownfield

- [ ] **Task 4: Implement recent projects list** (AC: #1, #2)
  - [ ] Store recent projects in Golang state
  - [ ] Add `project.getRecent` JSON-RPC method
  - [ ] Display with project name, path, last opened date
  - [ ] Allow removing projects from recent list

- [ ] **Task 5: Add project context editor** (AC: #3)
  - [ ] Create optional "Project Context" text area
  - [ ] Save context with project configuration
  - [ ] Show for brownfield projects

## Dev Notes

### UI Components (shadcn/ui)

```tsx
// src/renderer/screens/ProjectSelectScreen.tsx

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Textarea } from '@/components/ui/textarea';

export function ProjectSelectScreen() {
  const [recentProjects, setRecentProjects] = useState<RecentProject[]>([]);
  const [selectedProject, setSelectedProject] = useState<ProjectScanResult | null>(null);
  
  const handleSelectFolder = async () => {
    const path = await window.api.dialog.selectFolder();
    if (path) {
      const result = await window.api.project.scan(path);
      setSelectedProject(result);
      
      if (result.isBmad) {
        // Add to recent projects
        await window.api.project.addRecent(path);
        // Reload recent list
        const recent = await window.api.project.getRecent();
        setRecentProjects(recent);
      }
    }
  };
  
  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8">Select Project</h1>
      
      <div className="grid gap-6 md:grid-cols-2">
        {/* New Project Selection */}
        <Card>
          <CardHeader>
            <CardTitle>Open Project</CardTitle>
          </CardHeader>
          <CardContent>
            <Button onClick={handleSelectFolder} className="w-full">
              Select Project Folder
            </Button>
          </CardContent>
        </Card>
        
        {/* Recent Projects */}
        <Card>
          <CardHeader>
            <CardTitle>Recent Projects</CardTitle>
          </CardHeader>
          <CardContent>
            {recentProjects.length === 0 ? (
              <p className="text-muted-foreground">No recent projects</p>
            ) : (
              <ul className="space-y-2">
                {recentProjects.map((project) => (
                  <RecentProjectItem 
                    key={project.path} 
                    project={project}
                    onSelect={() => handleOpenRecent(project.path)}
                  />
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
      </div>
      
      {/* Project Info (shown after selection) */}
      {selectedProject && (
        <ProjectInfoCard project={selectedProject} />
      )}
    </div>
  );
}
```

### Project Info Component

```tsx
// src/renderer/components/ProjectInfoCard.tsx

export function ProjectInfoCard({ project }: { project: ProjectScanResult }) {
  const [context, setContext] = useState('');
  
  if (!project.isBmad) {
    return (
      <Card className="mt-6 border-destructive">
        <CardContent className="pt-6">
          <p className="text-destructive font-medium">Not a BMAD Project</p>
          <p className="text-muted-foreground mt-2">
            Run <code>bmad-init</code> to initialize this folder as a BMAD project.
          </p>
        </CardContent>
      </Card>
    );
  }
  
  return (
    <Card className="mt-6">
      <CardHeader>
        <CardTitle>Project Details</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Project Type Badge */}
        <div className="flex items-center gap-2">
          <span className="font-medium">Type:</span>
          <Badge variant={project.projectType === 'greenfield' ? 'default' : 'secondary'}>
            {project.projectType}
          </Badge>
        </div>
        
        {/* BMAD Version */}
        <div className="flex items-center gap-2">
          <span className="font-medium">BMAD Version:</span>
          <span>{project.bmadVersion}</span>
          {project.bmadCompatible ? (
            <CheckCircle className="h-4 w-4 text-green-500" />
          ) : (
            <XCircle className="h-4 w-4 text-red-500" />
          )}
        </div>
        
        {/* Existing Artifacts (Brownfield) */}
        {project.existingArtifacts && project.existingArtifacts.length > 0 && (
          <div>
            <span className="font-medium">Existing Artifacts:</span>
            <ul className="mt-2 space-y-1">
              {project.existingArtifacts.map((artifact) => (
                <li key={artifact.path} className="text-sm text-muted-foreground">
                  • {artifact.name} ({artifact.type})
                </li>
              ))}
            </ul>
          </div>
        )}
        
        {/* Project Context (Brownfield only) */}
        {project.projectType === 'brownfield' && (
          <div>
            <Label htmlFor="context">Project Context (optional)</Label>
            <Textarea
              id="context"
              placeholder="Describe your project context for better AI assistance..."
              value={context}
              onChange={(e) => setContext(e.target.value)}
              className="mt-2"
            />
            <Button 
              variant="outline" 
              size="sm" 
              className="mt-2"
              onClick={() => handleSaveContext(project.path, context)}
            >
              Save Context
            </Button>
          </div>
        )}
        
        {/* Continue Button */}
        <Button className="w-full mt-4">
          Continue to Dashboard
        </Button>
      </CardContent>
    </Card>
  );
}
```

### Preload API Additions

```typescript
// src/preload/index.ts (additions)

const api = {
  dialog: {
    selectFolder: (): Promise<string | null> => 
      ipcRenderer.invoke('dialog:selectFolder'),
  },
  project: {
    // ... existing
    getRecent: (): Promise<RecentProject[]> => 
      ipcRenderer.invoke('rpc:call', 'project.getRecent'),
    addRecent: (path: string): Promise<void> => 
      ipcRenderer.invoke('rpc:call', 'project.addRecent', { path }),
    removeRecent: (path: string): Promise<void> => 
      ipcRenderer.invoke('rpc:call', 'project.removeRecent', { path }),
    setContext: (path: string, context: string): Promise<void> => 
      ipcRenderer.invoke('rpc:call', 'project.setContext', { path, context }),
  },
};
```

### Main Process Dialog Handler

```typescript
// src/main/index.ts (additions)

import { dialog, ipcMain } from 'electron';

ipcMain.handle('dialog:selectFolder', async () => {
  const result = await dialog.showOpenDialog({
    properties: ['openDirectory'],
    title: 'Select BMAD Project Folder',
  });
  
  if (result.canceled || result.filePaths.length === 0) {
    return null;
  }
  
  return result.filePaths[0];
});
```

### Recent Projects Storage

```go
// internal/project/recent.go

type RecentProject struct {
    Path       string    `json:"path"`
    Name       string    `json:"name"`
    LastOpened time.Time `json:"lastOpened"`
    Context    string    `json:"context,omitempty"`
}

type RecentManager struct {
    configPath string
    projects   []RecentProject
    maxRecent  int
}

func (rm *RecentManager) Add(path string) error {
    // Remove if exists (to update timestamp)
    rm.Remove(path)
    
    // Add to front
    project := RecentProject{
        Path:       path,
        Name:       filepath.Base(path),
        LastOpened: time.Now(),
    }
    rm.projects = append([]RecentProject{project}, rm.projects...)
    
    // Trim to max
    if len(rm.projects) > rm.maxRecent {
        rm.projects = rm.projects[:rm.maxRecent]
    }
    
    return rm.save()
}
```

### File Structure

```
apps/desktop/src/renderer/
├── screens/
│   └── ProjectSelectScreen.tsx
├── components/
│   ├── ProjectInfoCard.tsx
│   └── RecentProjectItem.tsx
└── types/
    └── project.ts

apps/core/internal/project/
├── detector.go      # From Story 1.7
├── recent.go        # Recent projects management
└── context.go       # Project context storage
```

### Keyboard Navigation (NFR-U1)

- Tab through all interactive elements
- Enter to select folder or recent project
- Escape to cancel dialogs

### Testing Requirements

1. Test folder selection with native dialog
2. Test recent projects list updates
3. Test project info display for all types
4. Test context saving for brownfield
5. Test keyboard navigation

### Dependencies

- **Story 1.4, 1.6, 1.7**: Detection features must exist
- **Story 1.3**: IPC bridge must be working
- **Story 1.1**: shadcn/ui must be initialized

### References

- [prd.md#FR44] - User can select a project folder
- [prd.md#FR48] - User can describe project context manually
- [ux-design-specification.md] - UI patterns and components
- [prd.md#NFR-U1] - Keyboard navigation support

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Completion Notes List

- 

### Change Log

| Date | Change | Reason |
|------|--------|--------|
