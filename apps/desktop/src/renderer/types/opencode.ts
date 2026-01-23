/**
 * OpenCode Profile Types
 *
 * These types correspond to the Go structs in apps/core/internal/opencode/profiles.go
 *
 * SECURITY NOTE: The raw alias command is NOT exposed to prevent shell injection attacks.
 * Only the profile name (validated alphanumeric) is safe to use for execution.
 */

export interface OpenCodeProfile {
  name: string
  // alias field removed for security (shell injection prevention)
  available: boolean
  error?: string
  isDefault: boolean
}

export interface ProfilesResult {
  profiles: OpenCodeProfile[]
  defaultFound: boolean
  source: string
}
