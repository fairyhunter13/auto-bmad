import { ElectronAPI } from '@electron-toolkit/preload'

/**
 * Type-safe API surface exposed to the renderer via contextBridge.
 */
export interface Api {
  system: {
    ping: () => Promise<string>
    getVersion: () => Promise<{ version: string; commit: string; date: string }>
  }
  dialog: {
    selectFolder: () => Promise<string | null>
  }
  project: {
    detect: (path: string) => Promise<{ valid: boolean; errors?: string[] }>
    validate: (path: string) => Promise<{ valid: boolean; missing?: string[] }>
    detectDependencies: () => Promise<{
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
    }>
    scan: (path: string) => Promise<{
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
    }>
    getRecent: () => Promise<Array<{
      path: string
      name: string
      lastOpened: string
      context?: string
    }>>
    addRecent: (path: string) => Promise<void>
    removeRecent: (path: string) => Promise<void>
    setContext: (path: string, context: string) => Promise<void>
  }
  opencode: {
    getProfiles: () => Promise<{
      profiles: Array<{
        name: string
        alias: string
        available: boolean
        error?: string
        isDefault: boolean
      }>
      defaultFound: boolean
      source: string
    }>
    detect: () => Promise<{
      found: boolean
      version?: string
      path?: string
      compatible: boolean
      minVersion: string
      error?: string
    }>
  }
  network: {
    getStatus: () => Promise<{
      status: 'online' | 'offline' | 'checking'
      lastChecked: string
      latency?: number
    }>
  }
  settings: {
    get: () => Promise<{
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
    }>
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
    }>) => Promise<{
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
    }>
    reset: () => Promise<{
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
    }>
  }
  on: {
    backendCrash: (callback: (error: string) => void) => () => void
    backendStatus: (callback: (status: 'connected' | 'disconnected') => void) => () => void
    networkStatusChanged: (
      callback: (event: {
        previous: 'online' | 'offline' | 'checking'
        current: 'online' | 'offline' | 'checking'
      }) => void
    ) => () => void
  }
  invoke: <T = unknown>(method: string, params?: unknown) => Promise<T>
}

declare global {
  interface Window {
    electron: ElectronAPI
    api: Api
  }
}
