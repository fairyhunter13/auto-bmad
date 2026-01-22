/**
 * Type definitions for BMAD project structure detection.
 * These types match the JSON-RPC response from project.scan.
 */

/**
 * Project type classification.
 */
export type ProjectType = 'not-bmad' | 'greenfield' | 'brownfield'

/**
 * Artifact found in _bmad-output folder.
 */
export interface Artifact {
  /** Filename of the artifact */
  name: string

  /** Relative path from _bmad-output (e.g., "planning-artifacts/prd.md") */
  path: string

  /** Type of artifact (prd, architecture, epics, ux-design, product-brief, other) */
  type: string

  /** ISO 8601 timestamp of last modification (optional) */
  modified?: string
}

/**
 * Complete result from scanning a project directory.
 */
export interface ProjectScanResult {
  /** Whether the project has _bmad/ folder */
  isBmad: boolean

  /** Project type classification */
  projectType: ProjectType

  /** Detected BMAD version from manifest.yaml */
  bmadVersion?: string

  /** Whether BMAD version is compatible (6.0.0+) */
  bmadCompatible: boolean

  /** Minimum required BMAD version */
  minBmadVersion: string

  /** Project path that was scanned */
  path: string

  /** Whether _bmad/ folder exists */
  hasBmadFolder: boolean

  /** Whether _bmad-output/ folder exists */
  hasOutputFolder: boolean

  /** List of existing artifacts (only for brownfield projects) */
  existingArtifacts?: Artifact[]

  /** Error message if project is not BMAD or scan failed */
  error?: string
}
