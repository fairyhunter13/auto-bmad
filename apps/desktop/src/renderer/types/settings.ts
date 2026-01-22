/**
 * User settings for Auto-BMAD
 * Persisted to ~/.autobmad/_bmad-output/.autobmad/config.json
 */
export interface Settings {
  // Retry settings
  maxRetries: number // Default: 3
  retryDelay: number // Default: 5000 (ms)

  // Notification settings
  desktopNotifications: boolean // Default: true
  soundEnabled: boolean // Default: false

  // Timeout settings
  stepTimeoutDefault: number // Default: 300000 (5 min)
  heartbeatInterval: number // Default: 60000 (60s)

  // UI preferences
  theme: 'light' | 'dark' | 'system' // Default: "system"
  showDebugOutput: boolean // Default: false

  // Project memory
  lastProjectPath?: string
  projectProfiles: Record<string, string> // path -> profile name
  recentProjectsMax: number // Default: 10
}
