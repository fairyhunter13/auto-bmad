/**
 * Type definitions for dependency detection and validation.
 * These types match the JSON-RPC response from project.detectDependencies.
 */

/**
 * OpenCode CLI detection result.
 */
export interface OpenCodeDetection {
  /** Whether OpenCode CLI was found in PATH */
  found: boolean

  /** Detected version (e.g., "0.1.5") */
  version?: string

  /** Full path to the opencode executable */
  path?: string

  /** Whether the version meets minimum requirements */
  compatible: boolean

  /** Minimum required version */
  minVersion: string

  /** Error message if detection failed */
  error?: string
}

/**
 * Git detection result.
 */
export interface GitDetection {
  /** Whether Git was found in PATH */
  found: boolean

  /** Detected version (e.g., "2.39.0") */
  version?: string

  /** Full path to the git executable */
  path?: string

  /** Whether the version meets minimum requirements (2.0+) */
  compatible: boolean

  /** Minimum required version */
  minVersion: string

  /** Error message if detection failed */
  error?: string
}

/**
 * Complete dependency detection result from project.detectDependencies.
 */
export interface DependencyDetectionResult {
  /** OpenCode CLI detection result */
  opencode: OpenCodeDetection

  /** Git detection result */
  git: GitDetection
}
