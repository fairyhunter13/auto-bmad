import { ElectronAPI } from '@electron-toolkit/preload'

/**
 * Type-safe API surface exposed to the renderer via contextBridge.
 */
export interface Api {
  system: {
    ping: () => Promise<string>
    getVersion: () => Promise<{ version: string; commit: string; date: string }>
  }
  project: {
    detect: (path: string) => Promise<{ valid: boolean; errors?: string[] }>
    validate: (path: string) => Promise<{ valid: boolean; missing?: string[] }>
  }
  on: {
    backendCrash: (callback: (error: string) => void) => () => void
    backendStatus: (callback: (status: 'connected' | 'disconnected') => void) => () => void
  }
  invoke: <T = unknown>(method: string, params?: unknown) => Promise<T>
}

declare global {
  interface Window {
    electron: ElectronAPI
    api: Api
  }
}
