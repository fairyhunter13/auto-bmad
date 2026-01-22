/**
 * Scope Management Module
 *
 * Provides multi-scope parallel artifact system functionality
 * for isolated development workflows.
 *
 * @module scope
 *
 * ## Usage Contexts
 *
 * ### CLI Commands (tools/cli/commands/scope.js)
 * - ScopeValidator: ID and configuration validation
 * - ScopeManager: CRUD operations for scopes
 * - ScopeInitializer: Directory structure creation
 * - ScopeMigrator: Legacy artifact migration
 * - ScopeSync: Sync-up and sync-down operations
 * - EventLogger: Cross-scope event tracking
 *
 * ### IDE/Workflow Integration (agents, workflows)
 * - ScopeContext: Session-sticky scope resolution for LLM conversations
 * - ArtifactResolver: Read-any/write-own access control enforcement
 *
 * ### Internal/Utility
 * - StateLock: File locking for concurrent access safety
 *
 * ---
 *
 * ## Directory Structure Reference
 *
 * ### Base Structure
 * ```
 * project-root/
 * ├── _bmad/
 * │   ├── _config/
 * │   │   └── scopes.yaml          # Scope registry and configuration
 * │   └── _events/
 * │       ├── event-log.yaml       # Cross-scope event log
 * │       └── subscriptions.yaml   # Scope event subscriptions
 * ├── _bmad-output/
 * │   ├── _shared/                 # Shared artifacts layer (syncUp target)
 * │   │   └── {sourceScope}/       # Namespaced by promoting scope
 * │   │       └── {artifact}       # Promoted artifact
 * │   │       └── {artifact}.meta  # Promotion metadata
 * │   └── {scope}/                 # Individual scope directories
 * │       ├── .scope-meta.yaml     # Scope metadata
 * │       ├── .sync-meta.yaml      # Sync tracking metadata
 * │       ├── planning-artifacts/  # Planning documents
 * │       ├── implementation-artifacts/
 * │       ├── tests/
 * │       ├── shared/              # syncDown target directory
 * │       │   └── {sourceScope}/   # Pulled artifacts namespaced by source
 * │       └── project-context.md   # Scope-specific context
 * └── .bmad-scope                  # Session-sticky scope file
 * ```
 *
 * ### ScopeSync Directory Structures
 *
 * #### syncUp (Promote to Shared Layer)
 * Promotes artifacts from a scope to the shared layer for cross-scope access.
 *
 * **Source Path:**
 * `_bmad-output/{scopeId}/{relativePath}`
 *
 * **Target Path:**
 * `_bmad-output/_shared/{scopeId}/{relativePath}`
 *
 * **Metadata File:**
 * `_bmad-output/_shared/{scopeId}/{relativePath}.meta`
 *
 * **Example:**
 * ```
 * # Promoting architecture.md from 'auth' scope
 * Source: _bmad-output/auth/planning-artifacts/architecture.md
 * Target: _bmad-output/_shared/auth/planning-artifacts/architecture.md
 * Meta:   _bmad-output/_shared/auth/planning-artifacts/architecture.md.meta
 * ```
 *
 * The `.meta` file contains:
 * - source_scope: Original scope ID
 * - promoted_at: ISO timestamp
 * - original_path: Relative path within scope
 * - original_hash: MD5 hash for change detection
 * - version: Incremental version number
 *
 * #### syncDown (Pull from Shared Layer)
 * Pulls shared artifacts into a scope's local `shared/` directory.
 *
 * **Source Path:**
 * `_bmad-output/_shared/{sourceScope}/{relativePath}`
 *
 * **Target Path:**
 * `_bmad-output/{targetScope}/shared/{sourceScope}/{relativePath}`
 *
 * **Example:**
 * ```
 * # Pulling auth's architecture.md into 'payments' scope
 * Source: _bmad-output/_shared/auth/planning-artifacts/architecture.md
 * Target: _bmad-output/payments/shared/auth/planning-artifacts/architecture.md
 * ```
 *
 * The double-namespacing (`{targetScope}/shared/{sourceScope}/`) ensures:
 * - Clear provenance: Artifacts show which scope they came from
 * - No conflicts: Multiple scopes' artifacts can coexist
 * - Version tracking: `.sync-meta.yaml` tracks pulled versions
 *
 * #### Sync Metadata (.sync-meta.yaml)
 * Each scope maintains sync state in `.sync-meta.yaml`:
 * ```yaml
 * version: 1
 * lastSyncUp: "2025-01-22T10:00:00.000Z"
 * lastSyncDown: "2025-01-22T11:00:00.000Z"
 * promotedFiles:
 *   "planning-artifacts/architecture.md":
 *     promotedAt: "2025-01-22T10:00:00.000Z"
 *     hash: "abc123..."
 *     version: 2
 * pulledFiles:
 *   "auth/planning-artifacts/architecture.md":
 *     pulledAt: "2025-01-22T11:00:00.000Z"
 *     version: 2
 *     hash: "abc123..."
 * ```
 *
 * ---
 *
 * ## Concurrent Access Patterns
 *
 * ### StateLock Usage
 * All state file modifications use `StateLock` for safe concurrent access:
 *
 * ```javascript
 * const lock = new StateLock();
 *
 * // Basic locking
 * const result = await lock.withLock('/path/to/state.yaml', async () => {
 *   // Safe operations here - lock is held
 *   const data = await readFile();
 *   data.count++;
 *   await writeFile(data);
 *   return data;
 * }); // Lock automatically released
 *
 * // Optimistic versioning (for conflict detection)
 * const result = await lock.optimisticUpdate(
 *   '/path/to/state.yaml',
 *   expectedVersion,  // Version you read earlier
 *   newData           // Data to write
 * );
 * if (result.conflict) {
 *   // Another process modified the file - handle conflict
 * }
 *
 * // Automatic version management
 * const updated = await lock.updateYamlWithVersion('/path/to/state.yaml', (data) => {
 *   data.field = 'new value';
 *   return data;  // Version incremented automatically by writeYaml()
 * });
 * ```
 *
 * ### Lock Behavior
 * - **Stale Detection:** Locks older than 30 seconds are considered stale
 * - **Retry Strategy:** Exponential backoff (100ms to 1000ms) with 10 retries
 * - **Atomic Creation:** Uses `wx` flag for race-free lock acquisition
 * - **Crash Recovery:** Stale locks are automatically cleaned up
 *
 * ### Version Tracking
 * YAML state files include automatic version fields:
 * - `_version`: Incremented on each write (for optimistic concurrency)
 * - `_lastModified`: ISO timestamp of last modification
 *
 * ---
 *
 * ## Migration and Rollback
 *
 * ### ScopeMigrator Workflow
 * Migrates legacy non-scoped artifacts to the new scoped structure:
 *
 * ```javascript
 * const migrator = new ScopeMigrator({ projectRoot: '/path/to/project' });
 *
 * // Check if migration is needed
 * if (await migrator.needsMigration()) {
 *   // Analyze existing artifacts
 *   const analysis = await migrator.analyzeExisting();
 *   console.log(`Found ${analysis.files.length} files to migrate`);
 *
 *   // Perform migration (creates backup by default)
 *   const result = await migrator.migrate({
 *     scopeId: 'default',  // Target scope ID
 *     backup: true         // Create backup before migration
 *   });
 *
 *   console.log(`Migrated to: ${result.scopeId}`);
 *   console.log(`Backup at: ${result.backupPath}`);
 * }
 * ```
 *
 * ### Backup Structure
 * Backups are stored in `_bmad-output/_backup_migration_{timestamp}/`:
 * ```
 * _bmad-output/
 * └── _backup_migration_1705920000000/
 *     ├── planning-artifacts/      # Backed up legacy directories
 *     ├── implementation-artifacts/
 *     ├── tests/
 *     ├── project-context.md       # Backed up root-level files
 *     └── sprint-status.yaml
 * ```
 *
 * ### Rollback
 * If migration fails or needs to be undone:
 *
 * ```javascript
 * // Rollback using backup path from migration result
 * await migrator.rollback(result.backupPath);
 * ```
 *
 * Rollback behavior:
 * 1. Verifies backup exists at specified path
 * 2. Removes current versions of backed-up files/directories
 * 3. Restores all items from backup to original locations
 * 4. Removes backup directory after successful restore
 *
 * **Note:** Rollback does NOT remove the scope from `scopes.yaml`.
 * Manual cleanup may be needed if the scope was registered.
 */

const { ScopeValidator } = require('./scope-validator');
const { ScopeManager } = require('./scope-manager');
const { ScopeInitializer } = require('./scope-initializer');
const { ScopeMigrator } = require('./scope-migrator');
const { ScopeContext } = require('./scope-context');
const { ArtifactResolver } = require('./artifact-resolver');
const { StateLock } = require('./state-lock');
const { ScopeSync } = require('./scope-sync');
const { EventLogger } = require('./event-logger');

module.exports = {
  // Core CRUD operations
  ScopeValidator, // ID validation, config schema, circular dependency detection
  ScopeManager, // Create, read, update, delete scopes in scopes.yaml
  ScopeInitializer, // Create directory structure for scopes

  // Migration and sync
  ScopeMigrator, // Migrate legacy non-scoped artifacts to scoped structure
  ScopeSync, // Promote artifacts to shared layer, pull updates

  // IDE/Workflow integration (used by agents and workflow templates)
  ScopeContext, // Conversation-sticky scope resolution (parallel-safe)
  ArtifactResolver, // Enforce read-any/write-own access model

  // Event system
  EventLogger, // Log and track cross-scope events, subscriptions

  // Utilities
  StateLock, // File locking for safe concurrent access to state files
};
