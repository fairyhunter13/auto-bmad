/**
 * Preload Script - Secure IPC Bridge
 *
 * This script runs in a privileged context with access to Node.js APIs,
 * but exposes only a safe, type-safe API surface to the renderer process
 * via contextBridge.
 *
 * Security requirements (NFR-S5):
 * - contextIsolation: true (enforced in BrowserWindow)
 * - nodeIntegration: false (enforced in BrowserWindow)
 * - Only expose specific methods, never raw ipcRenderer
 */

import { contextBridge, ipcRenderer } from 'electron'

/**
 * Type-safe API surface exposed to the renderer.
 * All methods return Promises and communicate via IPC.
 */
const api = {
  /**
   * System methods for health checks and diagnostics
   */
  system: {
    /** Ping the backend to verify connectivity */
    ping: (): Promise<string> => ipcRenderer.invoke('rpc:call', 'system.ping'),

    /** Get backend version information */
    getVersion: (): Promise<{ version: string; commit: string; date: string }> =>
      ipcRenderer.invoke('rpc:call', 'system.version')
  },

  /**
   * Dialog methods for native file system dialogs
   */
  dialog: {
    /** Open native folder selection dialog */
    selectFolder: (): Promise<string | null> =>
      ipcRenderer.invoke('dialog:selectFolder')
  },

  /**
   * Project methods for BMAD project operations
   */
  project: {
    /** Detect BMAD project structure at the given path */
    detect: (path: string): Promise<{ valid: boolean; errors?: string[] }> =>
      ipcRenderer.invoke('rpc:call', 'project.detect', { path }),

    /** Validate project dependencies (Git, OpenCode) */
    validate: (path: string): Promise<{ valid: boolean; missing?: string[] }> =>
      ipcRenderer.invoke('rpc:call', 'project.validate', { path }),

    /** Detect and validate system dependencies (OpenCode, Git) */
    detectDependencies: (): Promise<{
      opencode: {
        found: boolean
        version?: string
        path?: string
        compatible: boolean
        minVersion: string
        error?: string
      }
      git: {
        found: boolean
        version?: string
        path?: string
        compatible: boolean
        minVersion: string
        error?: string
      }
    }> => ipcRenderer.invoke('rpc:call', 'project.detectDependencies'),

    /** Scan project directory for BMAD structure and artifacts */
    scan: (path: string): Promise<{
      isBmad: boolean
      projectType: 'not-bmad' | 'greenfield' | 'brownfield'
      bmadVersion?: string
      bmadCompatible: boolean
      minBmadVersion: string
      path: string
      hasBmadFolder: boolean
      hasOutputFolder: boolean
      existingArtifacts?: Array<{
        name: string
        path: string
        type: string
        modified?: string
      }>
      error?: string
    }> => ipcRenderer.invoke('rpc:call', 'project.scan', { path }),

    /** Get list of recent projects */
    getRecent: (): Promise<Array<{
      path: string
      name: string
      lastOpened: string
      context?: string
    }>> => ipcRenderer.invoke('rpc:call', 'project.getRecent'),

    /** Add project to recent list */
    addRecent: (path: string): Promise<void> =>
      ipcRenderer.invoke('rpc:call', 'project.addRecent', { path }),

    /** Remove project from recent list */
    removeRecent: (path: string): Promise<void> =>
      ipcRenderer.invoke('rpc:call', 'project.removeRecent', { path }),

    /** Set project context description */
    setContext: (path: string, context: string): Promise<void> =>
      ipcRenderer.invoke('rpc:call', 'project.setContext', { path, context })
  },

  /**
   * OpenCode methods for profile and CLI detection
   */
  opencode: {
    /** Get list of available OpenCode profiles from ~/.bash_aliases */
    getProfiles: (): Promise<{
      profiles: Array<{
        name: string
        alias: string
        available: boolean
        error?: string
        isDefault: boolean
      }>
      defaultFound: boolean
      source: string
    }> => ipcRenderer.invoke('rpc:call', 'opencode.getProfiles'),

    /** Detect OpenCode CLI installation */
    detect: (): Promise<{
      found: boolean
      version?: string
      path?: string
      compatible: boolean
      minVersion: string
      error?: string
    }> => ipcRenderer.invoke('rpc:call', 'opencode.detect')
  },

  /**
   * Network methods for connectivity monitoring
   */
  network: {
    /** Get current network connectivity status */
    getStatus: (): Promise<{
      status: 'online' | 'offline' | 'checking'
      lastChecked: string
      latency?: number
    }> => ipcRenderer.invoke('rpc:call', 'network.getStatus')
  },

  /**
   * Settings methods for user preferences
   */
  settings: {
    /** Get current settings */
    get: (): Promise<{
      maxRetries: number
      retryDelay: number
      desktopNotifications: boolean
      soundEnabled: boolean
      stepTimeoutDefault: number
      heartbeatInterval: number
      theme: string
      showDebugOutput: boolean
      lastProjectPath?: string
      projectProfiles: Record<string, string>
      recentProjectsMax: number
    }> => ipcRenderer.invoke('rpc:call', 'settings.get'),

    /** Update settings */
    set: (updates: Partial<{
      maxRetries: number
      retryDelay: number
      desktopNotifications: boolean
      soundEnabled: boolean
      stepTimeoutDefault: number
      heartbeatInterval: number
      theme: string
      showDebugOutput: boolean
      lastProjectPath: string
      projectProfiles: Record<string, string>
      recentProjectsMax: number
    }>): Promise<{
      maxRetries: number
      retryDelay: number
      desktopNotifications: boolean
      soundEnabled: boolean
      stepTimeoutDefault: number
      heartbeatInterval: number
      theme: string
      showDebugOutput: boolean
      lastProjectPath?: string
      projectProfiles: Record<string, string>
      recentProjectsMax: number
    }> => ipcRenderer.invoke('rpc:call', 'settings.set', updates),

    /** Reset all settings to defaults */
    reset: (): Promise<{
      maxRetries: number
      retryDelay: number
      desktopNotifications: boolean
      soundEnabled: boolean
      stepTimeoutDefault: number
      heartbeatInterval: number
      theme: string
      showDebugOutput: boolean
      lastProjectPath?: string
      projectProfiles: Record<string, string>
      recentProjectsMax: number
    }> => ipcRenderer.invoke('rpc:call', 'settings.reset')
  },

  /**
   * Event subscriptions for backend events.
   * Returns an unsubscribe function.
   */
  on: {
    /** Subscribe to backend crash events */
    backendCrash: (callback: (error: string) => void): (() => void) => {
      const handler = (_event: Electron.IpcRendererEvent, error: string): void => callback(error)
      ipcRenderer.on('backend:crash', handler)
      return () => ipcRenderer.removeListener('backend:crash', handler)
    },

    /** Subscribe to backend connection status changes */
    backendStatus: (callback: (status: 'connected' | 'disconnected') => void): (() => void) => {
      const handler = (
        _event: Electron.IpcRendererEvent,
        status: 'connected' | 'disconnected'
      ): void => callback(status)
      ipcRenderer.on('backend:status', handler)
      return () => ipcRenderer.removeListener('backend:status', handler)
    },

    /** Subscribe to network status changes */
    networkStatusChanged: (
      callback: (event: {
        previous: 'online' | 'offline' | 'checking'
        current: 'online' | 'offline' | 'checking'
      }) => void
    ): (() => void) => {
      const handler = (
        _event: Electron.IpcRendererEvent,
        event: { previous: string; current: string }
      ): void => callback(event as { previous: 'online' | 'offline' | 'checking'; current: 'online' | 'offline' | 'checking' })
      ipcRenderer.on('event:network.statusChanged', handler)
      return () => ipcRenderer.removeListener('event:network.statusChanged', handler)
    }
  },

  /**
   * Generic RPC call method for advanced use cases.
   * Prefer using the typed methods above when available.
   */
  invoke: <T = unknown>(method: string, params?: unknown): Promise<T> =>
    ipcRenderer.invoke('rpc:call', method, params) as Promise<T>
}

/**
 * Type declaration for the API exposed to the renderer.
 * Import this type in renderer code for type safety.
 */
export type Api = typeof api

// Expose API to renderer via contextBridge
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('api', api)
  } catch (error) {
    console.error('[Preload] Failed to expose API:', error)
  }
} else {
  // Fallback for non-isolated context (should not happen in production)
  console.warn('[Preload] Context isolation is disabled - this is a security risk')
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(window as any).api = api
}
