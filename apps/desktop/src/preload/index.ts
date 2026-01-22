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
   * Project methods for BMAD project operations
   */
  project: {
    /** Detect BMAD project structure at the given path */
    detect: (path: string): Promise<{ valid: boolean; errors?: string[] }> =>
      ipcRenderer.invoke('rpc:call', 'project.detect', { path }),

    /** Validate project dependencies (Git, OpenCode) */
    validate: (path: string): Promise<{ valid: boolean; missing?: string[] }> =>
      ipcRenderer.invoke('rpc:call', 'project.validate', { path })
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
