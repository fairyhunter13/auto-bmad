/**
 * OpenCode Profile Types
 *
 * These types correspond to the Go structs in apps/core/internal/opencode/profiles.go
 */

export interface OpenCodeProfile {
  name: string
  alias: string
  available: boolean
  error?: string
  isDefault: boolean
}

export interface ProfilesResult {
  profiles: OpenCodeProfile[]
  defaultFound: boolean
  source: string
}
